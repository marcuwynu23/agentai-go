package model

// ExampleResource represents an example domain entity.
type ExampleResource struct {
	Name string
}

// CreateExampleInput is the input for creating an example resource.
type CreateExampleInput struct {
	Name  string
	Force bool
}

// ListExampleInput is the input for listing example resources.
type ListExampleInput struct {
	Limit int
	All   bool
}

// ListExampleResult is the result of listing example resources.
type ListExampleResult struct {
	Items []ExampleResource
	Total int
}
