package identity

import (
	"log"
)

type contextKey string

const IdentityService contextKey = "IdentityService"
const IdentityTokenRefresh contextKey = "IdentityTokenRefresh"

type Service struct {
	TenantURL     string
	Client        *ClientWithResponses
	Logger        *log.Logger
	AuthnProvider AuthenticationProvider
}
