package identity

import (
	"context"
	"fmt"
	"io"
	"log"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "identity service context value " + k.name
}

var (
	ServiceKey = &contextKey{"IdentityService"}
	HeadersKey = &contextKey{"IdentityRequestHeaders"}
)

type Service struct {
	TenantURL     string
	TenantID      *string
	Client        *ClientWithResponses
	SessionToken  *string
	Logger        *log.Logger
	AuthnProvider AuthenticationProvider
}
type ServiceOption func(*Service) error

func NewService(ctx context.Context, idtenanturl string, opts ...ServiceOption) (*Service, error) {
	service := &Service{
		TenantURL: idtenanturl,
	}
	for _, o := range opts {
		if err := o(service); err != nil {
			return nil, err
		}
	}
	if service.Client == nil {
		clientWithResponses, errClient := NewClientWithResponses(idtenanturl)
		if errClient != nil {
			return nil, errClient
		}
		service.Client = clientWithResponses
	}

	// no logger set, so discard logs
	if service.Logger == nil {
		service.Logger = log.New(io.Discard, "", 0)
	}

	return service, nil
}

// ServiceWithClientWithResponses is a ServiceOption that sets the ClientWithResponses on the Service.
func ServiceWithClientWithResponses(client *ClientWithResponses) ServiceOption {
	return func(s *Service) error {
		s.Client = client
		return nil
	}
}

// ServiceWithLogger is a ServiceOption that sets the Logger on the Service.
func ServiceWithLogger(logger *log.Logger) ServiceOption {
	return func(s *Service) error {
		s.Logger = logger
		return nil
	}
}

// ServiceWithAuthnProvider is a ServiceOption that sets the AuthenticationProvider on the Service.
func ServiceWithAuthnProvider(provider AuthenticationProvider) ServiceOption {
	return func(s *Service) error {
		s.AuthnProvider = provider
		return nil
	}
}

// CreateRole creates a new role in the identity service. <https://api-docs.cyberark.com/docs/identity-api-reference/role-management/operations/create-a-role-store-role>
func (s *Service) CreateRole(ctx context.Context, reqCreateRole *PostRolesStoreRoleJSONRequestBody) (*RolesStoreRole, error) {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}
	ctx = context.WithValue(ctx, HeadersKey, headers)
	roleResp, err := s.Client.PostRolesStoreRoleWithResponse(ctx, *reqCreateRole,
		AddRequestHeaders,
		s.AuthnProvider.AuthenticateRequest,
	)
	if err != nil {
		return nil, err
	}
	return roleResp.JSON200, ReturnErrorWhenBodySuccessIsFalse(roleResp.JSON200)
}

// ReturnErrorWhenBodySuccessIsFalse returns an error if the body's Success field is false (even though HTTP status code is 200 OK).
func ReturnErrorWhenBodySuccessIsFalse(body *RolesStoreRole) error {
	if body == nil {
		return fmt.Errorf("failed to create role: body is nil")
	}
	success := body.Success
	if *success {
		return nil
	}
	Message := body.Message
	if Message != nil {
		return fmt.Errorf("failed to create role: %s", *Message)
	}
	MessageID := body.MessageID
	if MessageID != nil {
		return fmt.Errorf("failed to create role: %s", *MessageID)
	}
	ErrorCode := body.ErrorCode
	if ErrorCode != nil {
		return fmt.Errorf("failed to create role (no message available) error code: %s", *ErrorCode)
	}
	ErrorID := body.ErrorID
	if ErrorID != nil {
		return fmt.Errorf("failed to create role (no message available) error id: %s", *ErrorID)
	}
	return fmt.Errorf("failed to create role: unknown error")
}
