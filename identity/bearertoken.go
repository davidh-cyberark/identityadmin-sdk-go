package identity

import (
	"fmt"
	"time"
)

const BearerTokenDefaultExpirationSeconds int = 300 // 5 minutes

type BearerToken struct {
	Token     string
	TokenType string
	Expires   *time.Time
}

func NewBearerToken(token string, tokentype string, expireseconds int) *BearerToken {
	tok := BearerToken{
		Token:     token,
		TokenType: tokentype,
	}
	tok.SetExpiration(expireseconds)
	return &tok
}

func (token *BearerToken) SetExpiration(seconds int) {
	newtime := time.Now().Add(time.Duration(seconds) * time.Second)
	token.Expires = &newtime
}

func (token *BearerToken) String() string {
	return fmt.Sprintf("%s %s", token.TokenType, token.Token)
}

func (token *BearerToken) IsValid() bool {
	if token == nil || token.Token == "" || token.Expires == nil {
		return false
	}
	return token.Expires.After(time.Now())
}
