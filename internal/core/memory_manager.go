package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MemoryManager handles .memory.json and conversation history.
type MemoryManager struct {
	WorkspacePath string
	LogsPath      string
	memoryPath    string
}

// NewMemoryManager creates a MemoryManager for the given workspace.
func NewMemoryManager(workspacePath, logsPath string) *MemoryManager {
	return &MemoryManager{
		WorkspacePath: workspacePath,
		LogsPath:      logsPath,
		memoryPath:    filepath.Join(workspacePath, ".memory.json"),
	}
}

func (m *MemoryManager) ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// LoadMemory reads .memory.json from disk.
func (m *MemoryManager) LoadMemory() (map[string]interface{}, error) {
	data, err := os.ReadFile(m.memoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return m.emptyMemory(), nil
		}
		return nil, err
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		return m.emptyMemory(), nil
	}
	return out, nil
}

// SaveMemory writes memory to .memory.json.
func (m *MemoryManager) SaveMemory(memory map[string]interface{}) error {
	if memory == nil {
		var err error
		memory, err = m.LoadMemory()
		if err != nil {
			return err
		}
	}
	if err := m.ensureDir(m.WorkspacePath); err != nil {
		return err
	}
	data, err := json.MarshalIndent(memory, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.memoryPath, data, 0644)
}

// UpdateMemory appends an execution update to history.
func (m *MemoryManager) UpdateMemory(update map[string]interface{}) error {
	mem, err := m.LoadMemory()
	if err != nil {
		return err
	}
	hist, _ := mem["history"].([]interface{})
	if hist == nil {
		hist = []interface{}{}
	}
	update["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	hist = append(hist, update)
	if len(hist) > 100 {
		hist = hist[len(hist)-100:]
	}
	mem["history"] = hist
	mem["lastUpdated"] = time.Now().UTC().Format(time.RFC3339)
	return m.SaveMemory(mem)
}

// SetProjectName sets projectName in memory.
func (m *MemoryManager) SetProjectName(projectName string) error {
	mem, err := m.LoadMemory()
	if err != nil {
		return err
	}
	mem["projectName"] = projectName
	mem["lastUpdated"] = time.Now().UTC().Format(time.RFC3339)
	return m.SaveMemory(mem)
}

// AddConversationMessage appends a user or assistant message.
func (m *MemoryManager) AddConversationMessage(role, content string) error {
	mem, err := m.LoadMemory()
	if err != nil {
		return err
	}
	conv, _ := mem["conversationHistory"].([]interface{})
	if conv == nil {
		conv = []interface{}{}
	}
	conv = append(conv, map[string]interface{}{
		"role":      role,
		"content":   content,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
	if len(conv) > 20 {
		conv = conv[len(conv)-20:]
	}
	mem["conversationHistory"] = conv
	mem["lastUpdated"] = time.Now().UTC().Format(time.RFC3339)
	return m.SaveMemory(mem)
}

// GetConversationContext returns the last 10 messages as a string for AI context.
func (m *MemoryManager) GetConversationContext() (string, error) {
	mem, err := m.LoadMemory()
	if err != nil {
		return "No previous conversation.", nil
	}
	conv, _ := mem["conversationHistory"].([]interface{})
	if len(conv) == 0 {
		return "No previous conversation.", nil
	}
	start := len(conv) - 10
	if start < 0 {
		start = 0
	}
	var parts []string
	for i := start; i < len(conv); i++ {
		msg, _ := conv[i].(map[string]interface{})
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)
		parts = append(parts, role+": "+content)
	}
	return strings.Join(parts, "\n"), nil
}

func (m *MemoryManager) emptyMemory() map[string]interface{} {
	now := time.Now().UTC().Format(time.RFC3339)
	return map[string]interface{}{
		"history":             []interface{}{},
		"conversationHistory": []interface{}{},
		"projectState":        map[string]interface{}{},
		"projectName":         nil,
		"lastUpdated":         now,
		"createdAt":           now,
	}
}

// LogAction is a no-op in Go unless LogsPath is set (optional file logging can be added later).
func (m *MemoryManager) LogAction(category string, data map[string]interface{}) {
	_ = category
	_ = data
}
