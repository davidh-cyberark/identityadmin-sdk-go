package main

import (
	"testing"
	"time"

	"github.com/davidh-cyberark/identityadmin-sdk-go/identity"
)

func TestIsTokenValid(t *testing.T) {
	uc := identity.UserCredentialsAuthenticationProvider{}
	if uc.IsTokenValid() {
		t.Fatal("Expected token to be invalid")
	}

	uc.Token = &identity.BearerToken{
		Token: "abc",
	}
	if uc.IsTokenValid() {
		t.Fatal("Expected token to be invalid")
	}

	past := time.Now().Add(-5 * time.Minute)
	uc.Token.Expires = &past
	if uc.IsTokenValid() {
		t.Fatal("Expected token to be invalid")
	}

	fiveminutes := time.Now().Add(5 * time.Minute)
	uc.Token.Expires = &fiveminutes
	if !uc.IsTokenValid() {
		t.Fatal("Expected token to be valid")
	}
}
