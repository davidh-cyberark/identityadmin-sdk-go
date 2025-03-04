package identity

import (
	"context"
	"fmt"
	"net/http"
)

type UserCredentialsProvider struct {
	User string
	Pass string
}

// RefreshToken refreshes the session token for the user
func (uc *UserCredentialsProvider) RefreshToken(ctx context.Context) error {
	service := ctx.Value(ServiceKey).(*Service)
	logger := service.Logger
	client := service.Client

	reqStartAuth := PostSecurityStartAuthenticationJSONRequestBody{
		TenantId: &service.TenantURL,
		User:     uc.User,
		Version:  "1.0",
	}
	respAuthStart, respAuthStartErr := client.PostSecurityStartAuthenticationWithResponse(ctx,
		reqStartAuth)
	if respAuthStartErr != nil {
		return fmt.Errorf("failed call to AuthStart: %s", respAuthStartErr)
	}
	headers := headersToString(respAuthStart.HTTPResponse.Header)
	logger.Println("Start Authentication")
	logger.Printf("Headers: %s\n", headers)
	logger.Printf("Response: Status=%s\nBody=\n%s\n\n", respAuthStart.Status(), string(respAuthStart.Body))
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
		return fmt.Errorf("failed call to AdvanceAuth: %s", respAdvanceAuthErr)
	}
	if respAdvanceAuth.JSON200.Result.Token == nil {
		return fmt.Errorf("no token returned")
	}
	copy := *respAdvanceAuth.JSON200.Result.Token // make a copy of the string
	service.SessionToken = &copy
	return nil
}

// GetToken returns the session token for the user
func (uc *UserCredentialsProvider) GetToken(ctx context.Context) (string, error) {
	service := ctx.Value(ServiceKey).(*Service)
	if service.SessionToken != nil && *service.SessionToken != "" {
		return *service.SessionToken, nil
	}
	err := uc.RefreshToken(ctx)
	if err != nil {
		return "", err
	}
	if service.SessionToken != nil && *service.SessionToken != "" {
		return *service.SessionToken, nil
	}

	return "", nil
}

// UpdateRequestWithToken updates the request with the session token
func (uc *UserCredentialsProvider) UpdateRequestWithToken(ctx context.Context, req *http.Request) error {
	tok, err := uc.GetToken(ctx)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	return nil
}
