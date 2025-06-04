package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/davidh-cyberark/identityadmin-sdk-go/identity"
)

var (
	version string = "dev"
)

func main() {
	idtenanturl := flag.String("idtenanturl", "", "Identity URL")
	iduser := flag.String("iduser", "", "Identity user id")
	idpass := flag.String("idpass", "", "Identity user password")

	ver := flag.Bool("version", false, "Print version")
	debug := flag.Bool("d", false, "Enable debug settings")
	flag.Parse()

	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	// logger
	logger := log.New(os.Stderr, "[identity-client] ", log.LstdFlags)

	if !*debug {
		logger.SetOutput(io.Discard)
	}

	ctx := context.Background()

	// Create the Authentication provider
	userAuth := &identity.UserCredentialsAuthenticationProvider{
		User: *iduser,
		Pass: *idpass,
	}

	// Create the Identity client with the authentication provider
	client, clientErr := identity.NewClientWithResponses(*idtenanturl,
		identity.WithRequestEditorFn(userAuth.Intercept))
	if clientErr != nil {
		logger.Fatalf("failed to create client: %v", clientErr)
	}

	// Create the Identity service with the client and authentication provider
	service := &identity.Service{
		TenantURL:     *idtenanturl,
		Client:        client,
		Logger:        logger,
		AuthnProvider: userAuth,
	}

	ctx = context.WithValue(ctx, identity.IdentityService, service)

	// only when calling the refresh token, set the refreshtoken context
	refreshctx := context.WithValue(ctx, identity.IdentityTokenRefresh, "RefreshToken")

	err := identity.RefreshTokenWithUserCredentials(refreshctx, service, userAuth)
	if err != nil {
		log.Fatalf("failed to refresh token: %v", err)
	}

	fmt.Printf("Authorization: %s %s\n", userAuth.Token.TokenType, userAuth.Token.Token)
}
