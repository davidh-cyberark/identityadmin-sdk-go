package identity

import (
	"context"
	"fmt"
	"net/http"
)

var (
	ErrNoAuthnProvider = fmt.Errorf("no authentication provider")
)

type AuthenticationProvider interface {
	RefreshToken(ctx context.Context) error
	GetToken(ctx context.Context) (string, error)
	UpdateRequestWithToken(ctx context.Context, req *http.Request) error
}
