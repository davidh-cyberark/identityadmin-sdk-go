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
	AuthenticateRequest(ctx context.Context, req *http.Request) error
}
