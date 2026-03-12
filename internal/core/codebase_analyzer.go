package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// CodebaseAnalyzer scans and analyzes project files.
type CodebaseAnalyzer struct {
	WorkspacePath       string
	SupportedExtensions map[string]bool
	IgnorePatterns      []string
}

// NewCodebaseAnalyzer creates a new CodebaseAnalyzer.
func NewCodebaseAnalyzer(workspacePath string) *CodebaseAnalyzer {
	exts := map[string]bool{
		".js": true, ".ts": true, ".jsx": true, ".tsx": true,
		".json": true, ".py": true, ".java": true, ".go": true,
		".rs": true, ".cpp": true, ".c": true, ".h": true, ".hpp": true,
	}
	return &CodebaseAnalyzer{
		WorkspacePath:       workspacePath,
		SupportedExtensions: exts,
		IgnorePatterns: []string{
			"node_modules", ".git", ".next", "dist", "build",
			".cache", "coverage", ".env", ".memory.json",
		},
	}
}

func (c *CodebaseAnalyzer) shouldIgnore(path string) bool {
	name := filepath.Base(path)
	for _, p := range c.IgnorePatterns {
		if strings.Contains(name, p) {
			return true
		}
	}
	return false
}

// ScannedFile holds path and content for a scanned file.
type ScannedFile struct {
	Path      string
	FullPath  string
	Content   string
	Extension string
	Size      int64
}

// ScanCodebase returns all supported files under WorkspacePath.
func (c *CodebaseAnalyzer) ScanCodebase() ([]ScannedFile, error) {
	var files []ScannedFile
	if c.WorkspacePath == "" {
		return files, nil
	}
	info, err := os.Stat(c.WorkspacePath)
	if err != nil || !info.IsDir() {
		return files, nil
	}
	err = filepath.Walk(c.WorkspacePath, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if c.shouldIgnore(fullPath) {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if !c.SupportedExtensions[ext] {
			return nil
		}
		if c.shouldIgnore(fullPath) {
			return nil
		}
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(c.WorkspacePath, fullPath)
		rel = filepath.ToSlash(rel)
		files = append(files, ScannedFile{
			Path:      rel,
			FullPath:  fullPath,
			Content:   string(content),
			Extension: ext,
			Size:      info.Size(),
		})
		return nil
	})
	return files, err
}

// AnalysisResult is the result of Analyze().
type AnalysisResult struct {
	Files    []FileSummary `json:"files"`
	Structure AnalysisStructure `json:"structure"`
	Issues   []Issue        `json:"issues"`
	Summary  AnalysisSummary `json:"summary"`
}

// FileSummary is a short file summary.
type FileSummary struct {
	Path      string `json:"path"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Lines     int    `json:"lines"`
}

// AnalysisStructure holds structure info.
type AnalysisStructure struct {
	TotalFiles    int      `json:"totalFiles"`
	FileTypes     map[string]int `json:"fileTypes"`
	Directories   []string `json:"directories"`
	HasPackageJson bool    `json:"hasPackageJson"`
	HasReadme     bool    `json:"hasReadme"`
	HasTests      bool    `json:"hasTests"`
	Dependencies  []string `json:"dependencies"`
	EntryPoints   []string `json:"entryPoints"`
}

// Issue is a detected issue.
type Issue struct {
	Type     string `json:"type"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// AnalysisSummary summarizes counts.
type AnalysisSummary struct {
	TotalFiles     int `json:"totalFiles"`
	TotalIssues    int `json:"totalIssues"`
	CriticalIssues int `json:"criticalIssues"`
	Warnings       int `json:"warnings"`
}

// Analyze runs scan + structure + bug detection.
func (c *CodebaseAnalyzer) Analyze() (*AnalysisResult, error) {
	files, err := c.ScanCodebase()
	if err != nil {
		return nil, err
	}
	structure := c.analyzeStructure(files)
	bugs := c.detectBugs(files)
	sum := AnalysisSummary{
		TotalFiles:  len(files),
		TotalIssues: len(bugs),
	}
	for _, i := range bugs {
		if i.Severity == "high" {
			sum.CriticalIssues++
		} else {
			sum.Warnings++
		}
	}
	fileSummaries := make([]FileSummary, 0, len(files))
	for _, f := range files {
		lines := 0
		if f.Content != "" {
			lines = strings.Count(f.Content, "\n") + 1
		}
		fileSummaries = append(fileSummaries, FileSummary{
			Path: f.Path, Extension: f.Extension, Size: f.Size, Lines: lines,
		})
	}
	return &AnalysisResult{
		Files:     fileSummaries,
		Structure: structure,
		Issues:    bugs,
		Summary:   sum,
	}, nil
}

func (c *CodebaseAnalyzer) analyzeStructure(files []ScannedFile) AnalysisStructure {
	s := AnalysisStructure{
		FileTypes:   make(map[string]int),
		Directories: []string{},
	}
	for _, f := range files {
		s.FileTypes[f.Extension]++
		if f.Path == "package.json" {
			s.HasPackageJson = true
			var pkg struct {
				Dependencies map[string]string `json:"dependencies"`
				Main         string           `json:"main"`
				Bin          interface{}      `json:"bin"`
			}
			if json.Unmarshal([]byte(f.Content), &pkg) == nil {
				for k := range pkg.Dependencies {
					s.Dependencies = append(s.Dependencies, k)
				}
				if pkg.Main != "" {
					s.EntryPoints = append(s.EntryPoints, pkg.Main)
				}
			}
		}
		if strings.Contains(strings.ToLower(f.Path), "readme") {
			s.HasReadme = true
		}
		if strings.Contains(f.Path, "test") || strings.Contains(f.Path, "spec") {
			s.HasTests = true
		}
		dir := filepath.Dir(f.Path)
		if dir != "." {
			dir = filepath.ToSlash(dir)
			found := false
			for _, d := range s.Directories {
				if d == dir {
					found = true
					break
				}
			}
			if !found {
				s.Directories = append(s.Directories, dir)
			}
		}
	}
	return s
}

func (c *CodebaseAnalyzer) detectBugs(files []ScannedFile) []Issue {
	var issues []Issue
	for _, f := range files {
		lines := strings.Split(f.Content, "\n")
		for i, line := range lines {
			lineNum := i + 1
			if strings.Contains(line, "console.log") && !strings.Contains(f.Path, "test") {
				issues = append(issues, Issue{
					Type: "warning", File: f.Path, Line: lineNum,
					Message: "console.log found in non-test file", Severity: "low",
				})
			}
			if regexp.MustCompile(`//\s*(TODO|FIXME|HACK|XXX)`).MatchString(line) {
				issues = append(issues, Issue{
					Type: "todo", File: f.Path, Line: lineNum,
					Message: strings.TrimSpace(line), Severity: "info",
				})
			}
			if regexp.MustCompile(`catch\s*\([^)]*\)\s*\{\s*\}`).MatchString(line) {
				issues = append(issues, Issue{
					Type: "bug", File: f.Path, Line: lineNum,
					Message: "Empty catch block - errors are silently ignored", Severity: "medium",
				})
			}
		}
		if strings.Contains(f.Content, "await") && !strings.Contains(f.Content, "try") && !strings.Contains(f.Content, "catch") {
			issues = append(issues, Issue{
				Type: "warning", File: f.Path, Line: 0,
				Message: "Async code without try-catch error handling", Severity: "medium",
			})
		}
		if regexp.MustCompile(`localhost:\d+|127\.0\.0\.1|api-key|password\s*=\s*["']`).MatchString(strings.ToLower(f.Content)) {
			issues = append(issues, Issue{
				Type: "security", File: f.Path, Line: 0,
				Message: "Potential hardcoded credentials or localhost references", Severity: "high",
			})
		}
	}
	return issues
}

// FormatAnalysisForAI returns a string suitable for AI prompts.
func (c *CodebaseAnalyzer) FormatAnalysisForAI(a *AnalysisResult) string {
	var b strings.Builder
	b.WriteString("=== Codebase Analysis ===\n\n")
	b.WriteString("Total Files: " + strconv.Itoa(a.Summary.TotalFiles) + "\n")
	var ftParts []string
	for ext, count := range a.Structure.FileTypes {
		ftParts = append(ftParts, ext+": "+strconv.Itoa(count))
	}
	b.WriteString("File Types: " + strings.Join(ftParts, ", ") + "\n\n")
	if a.Structure.HasPackageJson {
		b.WriteString("Dependencies: " + strings.Join(a.Structure.Dependencies, ", ") + "\n")
		b.WriteString("Entry Points: " + strings.Join(a.Structure.EntryPoints, ", ") + "\n\n")
	}
	if len(a.Structure.Directories) > 0 {
		b.WriteString("Directory Structure:\n")
		for _, d := range a.Structure.Directories {
			b.WriteString("  - " + d + "\n")
		}
		b.WriteString("\n")
	}
	if len(a.Files) > 0 {
		b.WriteString("Key Files:\n")
		max := 20
		if len(a.Files) < max {
			max = len(a.Files)
		}
		for i := 0; i < max; i++ {
			f := a.Files[i]
			b.WriteString("  - " + f.Path + " (" + strconv.Itoa(f.Lines) + " lines, " + f.Extension + ")\n")
		}
		b.WriteString("\n")
	}
	if len(a.Issues) > 0 {
		b.WriteString("Issues Found: " + strconv.Itoa(a.Summary.TotalIssues) + "\n")
		b.WriteString("  - Critical: " + strconv.Itoa(a.Summary.CriticalIssues) + "\n")
		b.WriteString("  - Warnings: " + strconv.Itoa(a.Summary.Warnings) + "\n\n")
	}
	b.WriteString("=== End Analysis ===\n")
	return b.String()
}
