package identity

type contextKey string

const IdentityService contextKey = "IdentityService"
const IdentityTokenRefresh contextKey = "IdentityTokenRefresh"

type Service struct {
	TenantURL     string
	Client        *ClientWithResponses
	AuthnProvider AuthenticationProvider
}
