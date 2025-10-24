package identity

import (
	"context"
	"fmt"
	"net/http"
)

type UserCredentialsAuthenticationProvider struct {
	User  string
	Pass  string
	Token *BearerToken
}

// GetToken returns the token object (aka a getter)
func (uc *UserCredentialsAuthenticationProvider) GetToken() *BearerToken {
	return uc.Token
}

// RefreshToken refreshes the session token using user credentials; side effect is that it updates uc.Token
func (uc *UserCredentialsAuthenticationProvider) RefreshToken(ctx context.Context, req *http.Request) (*BearerToken, error) {
	if len(uc.User) == 0 || len(uc.Pass) == 0 {
		return nil, fmt.Errorf("user credentials are not set")
	}

	service := ctx.Value(IdentityService).(*Service)
	if service == nil {
		return nil, fmt.Errorf("service not set")
	}

	client := service.Client
	if client == nil {
		return nil, fmt.Errorf("client not set")
	}

	// Security Start Authentication
	reqStartAuth := PostSecurityStartAuthenticationJSONRequestBody{
		TenantId: &service.TenantURL,
		User:     uc.User,
		Version:  "1.0",
	}
	respAuthStart, respAuthStartErr := client.PostSecurityStartAuthenticationWithResponse(ctx, reqStartAuth)
	if respAuthStartErr != nil {
		return nil, fmt.Errorf("failed call to start authn: %s", respAuthStartErr)
	}

	persistLogin := true
	if respAuthStart.JSON200.Result.Challenges == nil {
		return nil, fmt.Errorf("no challenges returned")
	}
	if (*respAuthStart.JSON200.Result.Challenges)[0].Mechanisms == nil {
		return nil, fmt.Errorf("no mechanisms returned")
	}
	mechanisms := (*respAuthStart.JSON200.Result.Challenges)[0].Mechanisms
	if len(*mechanisms) == 0 {
		return nil, fmt.Errorf("no mechanisms returned")
	}

	// Security Advance Authentication
	reqAdvanceAuth := PostSecurityAdvanceAuthenticationJSONRequestBody{
		TenantId: respAuthStart.JSON200.Result.TenantId,
		// Possible Actions: "Unknown", "Answer", "StartOOB", "Poll", "ForgotPassword", "RetryOOB"
		Action:          "Answer",
		PersistentLogin: &persistLogin,
		SessionId:       *respAuthStart.JSON200.Result.SessionId,
		MechanismId:     *(*mechanisms)[0].MechanismId,
		Answer:          &uc.Pass,
	}
	respAdvanceAuth, respAdvanceAuthErr := client.PostSecurityAdvanceAuthenticationWithResponse(ctx,
		reqAdvanceAuth)
	if respAdvanceAuthErr != nil {
		return nil, fmt.Errorf("failed call to advance authn: %s", respAdvanceAuthErr)
	}

	if respAdvanceAuth.JSON200.Result.Token == nil {
		return nil, fmt.Errorf("no token returned")
	}

	uc.Token = NewBearerToken(*respAdvanceAuth.JSON200.Result.Token, "Bearer", BearerTokenDefaultExpirationSeconds)
	return uc.Token, nil
}
