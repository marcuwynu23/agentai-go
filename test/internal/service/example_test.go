package service_test

import (
	"testing"
	"github.com/marcuwynu23/cli-go-project-template/internal/model"
	"github.com/marcuwynu23/cli-go-project-template/internal/service"
)

func TestExampleService_Create(t *testing.T) {
	svc := service.NewExampleService()
	res, err := svc.Create(model.CreateExampleInput{Name: "my-resource", Force: false})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if res.Name != "my-resource" {
		t.Errorf("Name = %q, want my-resource", res.Name)
	}
}

func TestExampleService_Create_InvalidInput(t *testing.T) {
	svc := service.NewExampleService()
	_, err := svc.Create(model.CreateExampleInput{Name: ""})
	if err != service.ErrInvalidInput {
		t.Errorf("err = %v, want ErrInvalidInput", err)
	}
}

func TestExampleService_List(t *testing.T) {
	svc := service.NewExampleService()
	result, err := svc.List(model.ListExampleInput{Limit: 10, All: false})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if result.Total != 10 {
		t.Errorf("Total = %d, want 10", result.Total)
	}
	if result.Items == nil {
		t.Error("Items should not be nil")
	}
}

func TestExampleService_List_All(t *testing.T) {
	svc := service.NewExampleService()
	result, err := svc.List(model.ListExampleInput{Limit: 5, All: true})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if result.Total != 100 {
		t.Errorf("Total = %d, want 100 when All=true", result.Total)
	}
}
