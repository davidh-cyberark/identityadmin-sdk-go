package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/davidh-cyberark/identityadmin-sdk-go/identity"
)

var (
	version string = "dev"

	// Command line flags
	idtenanturl *string
	iduser      *string
	idpass      *string
	ver         *bool
)

func init() {
	idtenanturl = flag.String("idtenanturl", "", "Identity URL")
	iduser = flag.String("iduser", "", "Identity user id")
	idpass = flag.String("idpass", "", "Identity user password")

	ver = flag.Bool("version", false, "Print version")
	flag.Parse()
}

func main() {
	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	// Create the Authentication provider
	userAuth := &identity.UserCredentialsAuthenticationProvider{
		User: *iduser,
		Pass: *idpass,
	}

	servicectx, service, serviceErr := identity.NewService(context.Background(), *idtenanturl, userAuth, nil)
	if serviceErr != nil {
		log.Fatalf("failed to create service: %s", serviceErr.Error())
	}

	// Use the service ctx here so that the refresh by the auth provider can find the service object
	tok, tokErr := service.AuthnProvider.RefreshToken(servicectx, nil)
	if tokErr != nil {
		slog.Error("failed to refresh token", slog.Any("error", tokErr.Error()))
		os.Exit(1)
	}

	if tok == nil {
		slog.Error("token is nil", slog.Any("error", "refresh token returned nil token"))
		os.Exit(1)
	}

	fmt.Printf("Authorization: %s\n", userAuth.Token)
}
