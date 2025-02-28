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

	DEBUG := *debug

	// logger
	logger := log.New(os.Stderr, "[identity-client] ", log.LstdFlags)

	if !DEBUG {
		logger.SetOutput(io.Discard)
	}
	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	ctx := context.Background()

	// Create a new identity service
	service, errService := identity.NewService(ctx, *idtenanturl, identity.ServiceWithLogger(logger))
	if errService != nil {
		log.Fatalf("failed to create identity service: %s", errService)
	}
	ctx = context.WithValue(ctx, identity.ServiceKey, service)

	// Implements the AuthenticationProvider interface
	authnProvider := identity.UserCredentialsProvider{
		User: *iduser,
		Pass: *idpass,
	}

	// Get Token returns the session token for the user.
	tok, tokErr := authnProvider.GetToken(ctx)
	if tokErr != nil {
		log.Fatalf("failed to get token: %s", tokErr)
	}
	fmt.Printf("%s\n", tok)
}
