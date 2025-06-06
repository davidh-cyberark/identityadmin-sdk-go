package identity

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type UserCredentialsAuthenticationProvider struct {
	User  string
	Pass  string
	Token *BearerToken
}

// Intercept modifies the request with the session token; implements RequestEditorFn type
func (uc *UserCredentialsAuthenticationProvider) Intercept(ctx context.Context, req *http.Request) error {
	// When refreshing token, no need to run through intercept logic, so, just return
	refresh := ctx.Value(IdentityTokenRefresh)
	if refresh != nil && refresh == "RefreshToken" {
		return nil
	}

	// check for invalid or expired token, and call refresh token
	if !uc.IsTokenValid() {
		service := ctx.Value(IdentityService).(*Service)
		if service != nil {
			refreshctx := context.WithValue(ctx, IdentityTokenRefresh, "RefreshToken")

			// side effect is that uc.Token is updated with fresh token
			err := RefreshTokenWithUserCredentials(refreshctx, service, uc)
			if err != nil {
				return fmt.Errorf("failed to refresh token: %w", err)
			}
		}
	}

	// If token is valid at this point, add the auth request header
	if uc.IsTokenValid() {
		req.Header.Set("Authorization", uc.Token.TokenType+" "+uc.Token.Token)
	}
	return nil
}

func (uc *UserCredentialsAuthenticationProvider) IsTokenValid() bool {
	if uc.Token == nil || uc.Token.Token == "" || uc.Token.Expires == nil {
		return false
	}
	return uc.Token.Expires.After(time.Now())
}

// RefreshToken refreshes the session token for the user
func RefreshTokenWithUserCredentials(ctx context.Context, service *Service, uc *UserCredentialsAuthenticationProvider) error {
	if len(uc.User) == 0 || len(uc.Pass) == 0 {
		return fmt.Errorf("user credentials are not set")
	}
	if service == nil {
		return fmt.Errorf("service not set")
	}

	logger := service.Logger
	if logger == nil {
		return fmt.Errorf("logger not set")
	}
	client := service.Client
	if client == nil {
		return fmt.Errorf("client not set")
	}

	// Security Start Authentication
	reqStartAuth := PostSecurityStartAuthenticationJSONRequestBody{
		TenantId: &service.TenantURL,
		User:     uc.User,
		Version:  "1.0",
	}
	respAuthStart, respAuthStartErr := client.PostSecurityStartAuthenticationWithResponse(ctx, reqStartAuth)
	if respAuthStartErr != nil {
		return fmt.Errorf("failed call to start authn: %s", respAuthStartErr)
	}

	headers := headersToString(respAuthStart.HTTPResponse.Header)
	logger.Println("Start Authentication")
	logger.Printf("--- HEADERS BEGIN ---\n%s\n--- HEADERS END ---", headers)
	logger.Printf("--- RESPONSE BEGIN ---\nStatus=%s\nBody=\n%s\n--- RESPONSE END ---\n", respAuthStart.Status(), string(respAuthStart.Body))

	persistLogin := true
	if respAuthStart.JSON200.Result.Challenges == nil {
		return fmt.Errorf("no challenges returned")
	}
	if (*respAuthStart.JSON200.Result.Challenges)[0].Mechanisms == nil {
		return fmt.Errorf("no mechanisms returned")
	}
	mechanisms := (*respAuthStart.JSON200.Result.Challenges)[0].Mechanisms
	if len(*mechanisms) == 0 {
		return fmt.Errorf("no mechanisms returned")
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
		return fmt.Errorf("failed call to advance authn: %s", respAdvanceAuthErr)
	}
	headers = headersToString(respAdvanceAuth.HTTPResponse.Header)
	logger.Println("Start Advance Authentication")
	logger.Printf("--- HEADERS BEGIN ---\n%s\n--- HEADERS END ---", headers)
	logger.Printf("--- RESPONSE BEGIN ---\nStatus=%s\nBody=\n%s\n--- RESPONSE END ---\n", respAdvanceAuth.Status(), string(respAdvanceAuth.Body))

	if respAdvanceAuth.JSON200.Result.Token == nil {
		return fmt.Errorf("no token returned")
	}

	fiveminutes := time.Now().Add(5 * time.Minute)
	tok := BearerToken{
		Token:     *respAdvanceAuth.JSON200.Result.Token,
		TokenType: "Bearer",
		Expires:   &fiveminutes,
	}

	uc.Token = &tok
	return nil
}
