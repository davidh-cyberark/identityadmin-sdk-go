package identity

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type contextKey string

const IdentityService contextKey = "IdentityService"
const IdentityTokenRefresh contextKey = "IdentityTokenRefresh"
const IdentityBearerToken contextKey = "IdentityBearerToken"

type Service struct {
	TenantURL     string
	Client        *ClientWithResponses
	Logger        *slog.Logger
	AuthnProvider AuthenticationProvider
}

func NewService(ctx context.Context, tenanturl string, authnprovider AuthenticationProvider, logger *slog.Logger) (context.Context, *Service, error) {
	client, clientErr := NewClientWithResponses(tenanturl, WithRequestEditorFn(AuthorizationIntercept))
	if clientErr != nil {
		return ctx, nil, fmt.Errorf("failed to create client: %v", clientErr)
	}

	// nil logger passed in, let's interpret that to mean that logs should be discarded
	if logger == nil {
		discardHandler := slog.NewJSONHandler(io.Discard, nil)

		// Create a new logger using the discard handler
		logger = slog.New(discardHandler)
	}

	// Create the Identity service with the client and authentication provider
	service := &Service{
		TenantURL:     tenanturl,
		Client:        client,
		Logger:        logger,
		AuthnProvider: authnprovider,
	}

	ctx = context.WithValue(ctx, IdentityService, service)

	return ctx, service, nil
}

func (srv *Service) GetValidToken() *BearerToken {
	authnprovider := srv.AuthnProvider
	token := authnprovider.GetToken()
	if token != nil && token.IsValid() {
		return token
	}

	// if srv.AuthnProvider
	// if !uc.IsTokenValid() {
	// 	service := ctx.Value(IdentityService).(*Service)
	// 	if service != nil {
	// 		refreshctx := context.WithValue(ctx, IdentityTokenRefresh, "RefreshToken")

	// 		// side effect is that uc.Token is updated with fresh token
	// 		err := RefreshTokenWithUserCredentials(refreshctx, service, uc)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to refresh token: %w", err)
	// 		}
	// 	}
	// }

	// // If token is valid at this point, add the auth request header
	// if uc.IsTokenValid() {
	// 	req.Header.Set("Authorization", uc.Token.TokenType+" "+uc.Token.Token)
	// }

	return token
}

// AuthorizationIntercept modifies the request with the session token; implements RequestEditorFn type
func AuthorizationIntercept(ctx context.Context, req *http.Request) error {
	service := ctx.Value(IdentityService).(*Service)

	service.GetValidToken()

	// check for invalid or expired token, and call refresh token
	return nil
}
