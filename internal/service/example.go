package service

import (
	"github.com/marcuwynu23/cli-go-project-template/internal/model"
)

// ExampleUseCase defines business operations for example resources.
type ExampleUseCase interface {
	Create(in model.CreateExampleInput) (*model.ExampleResource, error)
	List(in model.ListExampleInput) (*model.ListExampleResult, error)
}

// ExampleService implements ExampleUseCase (business logic only, no I/O).
type ExampleService struct{}

// NewExampleService returns a new ExampleUseCase implementation.
func NewExampleService() *ExampleService {
	return &ExampleService{}
}

// Create creates an example resource. In a real app this would persist.
func (s *ExampleService) Create(in model.CreateExampleInput) (*model.ExampleResource, error) {
	if in.Name == "" {
		return nil, ErrInvalidInput
	}
	// Placeholder: validate, check exists, persist. For template we just return the model.
	return &model.ExampleResource{Name: in.Name}, nil
}

// List returns example resources. In a real app this would query a store.
func (s *ExampleService) List(in model.ListExampleInput) (*model.ListExampleResult, error) {
	limit := in.Limit
	if in.All {
		limit = 100
	}
	// Placeholder: fetch from store. For template we return empty slice with count.
	items := make([]model.ExampleResource, 0, limit)
	return &model.ListExampleResult{Items: items, Total: limit}, nil
}
