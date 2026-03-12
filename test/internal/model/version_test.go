package model_test

import (
	"testing"
	"github.com/marcuwynu23/cli-go-project-template/internal/model"
)

func TestVersionInfo_ZeroValue(t *testing.T) {
	var info model.VersionInfo
	if info.Version != "" {
		t.Errorf("zero value Version should be empty")
	}
}
