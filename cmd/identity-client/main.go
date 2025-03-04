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

	newrolename := flag.String("newrolename", "", "New role name")
	newroledesc := flag.String("newroledesc", "", "New role description")

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

	// Implements the AuthenticationProvider interface
	var authnProvider identity.AuthenticationProvider = &identity.UserCredentialsProvider{
		User: *iduser,
		Pass: *idpass,
	}

	// Create a new identity service
	service, errService := identity.NewService(ctx, *idtenanturl, identity.ServiceWithLogger(logger), identity.ServiceWithAuthnProvider(authnProvider))
	if errService != nil {
		log.Fatalf("failed to create identity service: %s", errService)
	}
	ctx = context.WithValue(ctx, identity.ServiceKey, service)

	// Get Token returns the session token for the user.
	tok, tokErr := service.AuthnProvider.GetToken(ctx)
	if tokErr != nil {
		log.Fatalf("failed to get token: %s", tokErr)
	}
	fmt.Printf("%s\n", tok)

	createRoleErr := CreateRole(ctx, *newrolename, *newroledesc)
	if createRoleErr != nil {
		log.Fatalf("failed to create role: %s", createRoleErr)
	}
}

func CreateRole(ctx context.Context, name, desc string) error {
	service, ok := ctx.Value(identity.ServiceKey).(*identity.Service)
	if !ok {
		return fmt.Errorf("failed to get identity service")
	}

	reqCreateRole := identity.PostRolesStoreRoleJSONRequestBody{
		Name:        name,
		Description: &desc,
	}

	roleResp, roleErr := service.CreateRole(ctx, &reqCreateRole)
	if roleErr != nil {
		return roleErr
	}
	msg := "Role created:"
	success := roleResp.Success
	if success != nil {
		msg = fmt.Sprintf("%s success: %t", msg, *success)
	}
	result := roleResp.Result
	if result != nil && result.Rowkey != nil {
		msg = fmt.Sprintf("%s, rowkey: %s", msg, *result.Rowkey)
	}
	fmt.Printf("%s\n", msg)
	return nil
}
