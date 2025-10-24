package identity

import (
	"context"
	"net/http"
)

type AuthenticationProvider interface {
	// Refresh the bearer token
	RefreshToken(ctx context.Context, req *http.Request) (*BearerToken, error)

	// Getter
	GetToken() *BearerToken
}
