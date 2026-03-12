package mcp

import (
	"os"
	"path/filepath"
)

// FileResult is the result of a file operation.
type FileResult struct {
	Success  bool   `json:"success"`
	FilePath string `json:"filePath,omitempty"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

// FileServer handles file create/read/modify (in-process).
type FileServer struct {
	WorkspacePath string
}

// NewFileServer creates a FileServer using cwd as workspace.
func NewFileServer(workspacePath string) *FileServer {
	if workspacePath == "" {
		workspacePath, _ = os.Getwd()
	}
	return &FileServer{WorkspacePath: workspacePath}
}

func (f *FileServer) ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}

// HandleCreateFile creates a file at path with content.
func (f *FileServer) HandleCreateFile(path, content string) FileResult {
	fullPath := filepath.Join(f.WorkspacePath, path)
	if err := f.ensureDir(fullPath); err != nil {
		return FileResult{Success: false, Error: err.Error(), FilePath: path}
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return FileResult{Success: false, Error: err.Error(), FilePath: path}
	}
	return FileResult{Success: true, FilePath: fullPath, Message: "File created: " + path}
}

// HandleModifyFile overwrites or appends to an existing file.
func (f *FileServer) HandleModifyFile(path, content string, append bool) FileResult {
	fullPath := filepath.Join(f.WorkspacePath, path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return FileResult{Success: false, Error: "File does not exist: " + path, FilePath: path}
	}
	if append {
		existing, err := os.ReadFile(fullPath)
		if err != nil {
			return FileResult{Success: false, Error: err.Error(), FilePath: path}
		}
		content = string(existing) + "\n" + content
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return FileResult{Success: false, Error: err.Error(), FilePath: path}
	}
	return FileResult{Success: true, FilePath: fullPath, Message: "File modified: " + path}
}

// ReadFile returns file content or error.
func (f *FileServer) ReadFile(path string) (content string, ok bool) {
	fullPath := filepath.Join(f.WorkspacePath, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", false
	}
	return string(data), true
}
