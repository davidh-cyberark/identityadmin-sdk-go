package identity

import "context"

type AuthenticationProvider interface {
	RefreshToken(ctx context.Context) error
	GetToken(ctx context.Context) (string, error)
}
