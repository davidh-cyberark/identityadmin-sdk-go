package identity

import (
	"context"
	"net/http"
	"time"
)

type BearerToken struct {
	Token     string
	TokenType string
	Expires   *time.Time
}

type AuthenticationProvider interface {
	Intercept(ctx context.Context, req *http.Request) error
}
