package identity

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrNoAuthnProvider = fmt.Errorf("no authentication provider")
)

type BearerToken struct {
	Token     string
	TokenType string
	Expires   *time.Time
}

type AuthenticationProvider interface {
	Intercept(ctx context.Context, req *http.Request) error
}
