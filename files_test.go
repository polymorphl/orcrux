package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApp_SaveFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "files_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	app := &App{}

	tests := []struct {
		name    string
		path    string
		content []byte
		wantErr bool
	}{
		{
			name:    "valid file save",
			path:    filepath.Join(tempDir, "test.txt"),
			content: []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			content: []byte("content"),
			wantErr: true,
		},
		{
			name:    "empty content",
			path:    filepath.Join(tempDir, "empty.txt"),
			content: []byte{},
			wantErr: true,
		},
		{
			name:    "nil content",
			path:    filepath.Join(tempDir, "nil.txt"),
			content: nil,
			wantErr: true,
		},
		{
			name:    "special characters in content",
			path:    filepath.Join(tempDir, "special.txt"),
			content: []byte("Special chars: 🚀🌟💻\n\t\r"),
			wantErr: false,
		},
		{
			name:    "large content",
			path:    filepath.Join(tempDir, "large.txt"),
			content: make([]byte, 1024*1024), // 1MB
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.SaveFile(tt.path, tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect success, verify the file was created with correct content
			if !tt.wantErr {
				// Check if file exists
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("SaveFile() file was not created at %s", tt.path)
					return
				}

				// Check file content
				readContent, err := os.ReadFile(tt.path)
				if err != nil {
					t.Errorf("SaveFile() failed to read created file: %v", err)
					return
				}

				if string(readContent) != string(tt.content) {
					t.Errorf("SaveFile() content mismatch, got %s, want %s",
						string(readContent), string(tt.content))
				}

				// Check file permissions (should be 0644)
				info, err := os.Stat(tt.path)
				if err != nil {
					t.Errorf("SaveFile() failed to get file info: %v", err)
					return
				}

				expectedMode := os.FileMode(0644)
				if info.Mode().Perm() != expectedMode {
					t.Errorf("SaveFile() file permissions mismatch, got %v, want %v",
						info.Mode().Perm(), expectedMode)
				}
			}
		})
	}
}

func TestApp_SaveFile_Overwrite(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "files_test_overwrite")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	app := &App{}
	testPath := filepath.Join(tempDir, "overwrite.txt")
	initialContent := []byte("Initial content")
	newContent := []byte("New content")

	// First save
	err = app.SaveFile(testPath, initialContent)
	if err != nil {
		t.Fatalf("Failed to save initial file: %v", err)
	}

	// Overwrite with new content
	err = app.SaveFile(testPath, newContent)
	if err != nil {
		t.Fatalf("Failed to overwrite file: %v", err)
	}

	// Verify new content
	readContent, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read overwritten file: %v", err)
	}

	if string(readContent) != string(newContent) {
		t.Errorf("File overwrite failed, got %s, want %s",
			string(readContent), string(newContent))
	}
}

func TestApp_SaveFile_InvalidPath(t *testing.T) {
	app := &App{}

	// Try to save to a directory that doesn't exist
	invalidPath := "/nonexistent/directory/test.txt"
	err := app.SaveFile(invalidPath, []byte("test"))

	if err == nil {
		t.Error("SaveFile() should fail with invalid path")
	}
}

func TestApp_UploadFile_NotImplemented(t *testing.T) {
	// Note: UploadFile() cannot be fully tested in unit tests because it
	// requires user interaction and native file dialogs. This test documents
	// the limitation and ensures the function exists.

	app := &App{}

	// We can only test that the function exists and has the right signature
	// The actual file dialog functionality would need integration tests
	_ = app.UploadFile
}

func TestApp_SaveFile_EdgeCases(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "files_test_edge")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	app := &App{}

	tests := []struct {
		name    string
		path    string
		content []byte
		wantErr bool
	}{
		{
			name:    "zero byte content",
			path:    filepath.Join(tempDir, "zero.txt"),
			content: make([]byte, 0),
			wantErr: true,
		},
		{
			name:    "path with spaces",
			path:    filepath.Join(tempDir, "file with spaces.txt"),
			content: []byte("content"),
			wantErr: false,
		},
		{
			name:    "path with special characters",
			path:    filepath.Join(tempDir, "file-@#$%^&*().txt"),
			content: []byte("content"),
			wantErr: false,
		},
		{
			name:    "very long filename",
			path:    filepath.Join(tempDir, "a"+strings.Repeat("b", 200)+".txt"),
			content: []byte("content"),
			wantErr: false, // This might fail on some systems, but we test the behavior
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.SaveFile(tt.path, tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// If successful, verify file was created
			if !tt.wantErr && err == nil {
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("SaveFile() file was not created at %s", tt.path)
				}
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkSaveFile(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "files_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	app := &App{}
	testContent := []byte("Benchmark test content")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPath := filepath.Join(tempDir, fmt.Sprintf("test_%d.txt", i))
		err := app.SaveFile(benchmarkPath, testContent)
		if err != nil {
			b.Fatalf("SaveFile failed: %v", err)
		}
	}
}

func BenchmarkSaveFile_LargeContent(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "files_benchmark_large")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	app := &App{}
	testContent := make([]byte, 1024*1024) // 1MB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPath := filepath.Join(tempDir, fmt.Sprintf("large_%d.txt", i))
		err := app.SaveFile(benchmarkPath, testContent)
		if err != nil {
			b.Fatalf("SaveFile failed: %v", err)
		}
	}
}
