package cmd

import (
	"github.com/marcuwynu23/cli-go-project-template/internal/service"
	"github.com/marcuwynu23/cli-go-project-template/internal/view"
)

// Deps holds service and view dependencies for the CLI (injectable for testing).
type Deps struct {
	VersionProvider service.VersionProvider
	ExampleUseCase  service.ExampleUseCase
	VersionView     *view.VersionRenderer
	ExampleView     *view.ExampleRenderer
}

// defaultDeps is set in init(); tests may replace it.
var defaultDeps *Deps

func initDeps() {
	if defaultDeps != nil {
		return
	}
	defaultDeps = &Deps{
		VersionProvider: service.NewVersionService(Version, Commit, BuildDate),
		ExampleUseCase:  service.NewExampleService(),
		VersionView:     view.NewVersionRenderer(),
		ExampleView:     view.NewExampleRenderer(),
	}
}

func deps() *Deps {
	initDeps()
	return defaultDeps
}

// ResetDepsForTest sets defaultDeps to nil so the next deps() call recreates them (e.g. with test Version).
// Only use from tests (e.g. test/cmd).
func ResetDepsForTest() {
	defaultDeps = nil
}
