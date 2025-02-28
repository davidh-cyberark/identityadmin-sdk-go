package identity

import (
	"context"
	"io"
	"log"
)

type ServiceKeyType string

const (
	ServiceKey ServiceKeyType = "Service"
)

type Service struct {
	TenantURL    string
	TenantID     *string
	Client       *ClientWithResponses
	SessionToken string
	Logger       *log.Logger
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
		clientWithResponses, errClient := NewClientWithResponses(string(idtenanturl))
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
func ServiceWithClientWithResponses(client *ClientWithResponses) ServiceOption {
	return func(s *Service) error {
		s.Client = client
		return nil
	}
}
func ServiceWithLogger(logger *log.Logger) ServiceOption {
	return func(s *Service) error {
		s.Logger = logger
		return nil
	}
}
