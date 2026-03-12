package mcp

import (
	"path/filepath"
	"strings"
)

// TestServer creates test files (uses FileServer).
type TestServer struct {
	WorkspacePath string
	TestFramework string
	fileServer    *FileServer
}

// NewTestServer creates a TestServer.
func NewTestServer(workspacePath string) *TestServer {
	if workspacePath == "" {
		workspacePath = "."
	}
	return &TestServer{
		WorkspacePath: workspacePath,
		TestFramework: "jest",
		fileServer:    NewFileServer(workspacePath),
	}
}

// TestResult is the result of test creation.
type TestResult struct {
	Success   bool   `json:"success"`
	TestPath  string `json:"testPath,omitempty"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
	Framework string `json:"framework,omitempty"`
}

// HandleCreateTests creates a test file at testPath or derives from targetFile.
func (t *TestServer) HandleCreateTests(testPath, testContent, targetFile string) TestResult {
	finalPath := testPath
	if finalPath == "" && targetFile != "" {
		base := filepath.Base(targetFile)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		finalPath = filepath.Join("tests", name+".test.js")
		finalPath = filepath.ToSlash(finalPath)
	}
	if finalPath == "" {
		finalPath = "tests/test.js"
	}
	content := testContent
	if content == "" {
		content = t.generateTestTemplate(targetFile)
	}
	res := t.fileServer.HandleCreateFile(finalPath, content)
	return TestResult{
		Success:   res.Success,
		TestPath:  finalPath,
		Message:   res.Message,
		Error:     res.Error,
		Framework: t.TestFramework,
	}
}

func (t *TestServer) generateTestTemplate(targetFile string) string {
	if targetFile == "" {
		targetFile = "Component"
	}
	if t.TestFramework == "jest" {
		return "describe('" + targetFile + "', () => {\n  test('should work correctly', () => {\n    expect(true).toBe(true);\n  });\n});\n"
	}
	if t.TestFramework == "mocha" {
		return "const { expect } = require('chai');\n\ndescribe('" + targetFile + "', () => {\n  it('should work correctly', () => {\n    expect(true).to.be.true;\n  });\n});\n"
	}
	return "// Test file for " + targetFile + "\n// Add your tests here\n"
}
