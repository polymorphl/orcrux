package main

import (
	"errors"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SaveFile writes the provided content to a file at the specified path.
//
// This function creates a new file or overwrites an existing file with the given
// content. The file is created with standard permissions (0644) which provides
// read/write access for the owner and read access for group and others.
//
// Parameters:
//   - path: The file path where content will be written (cannot be empty)
//   - content: The byte array to write to the file (cannot be empty)
//
// Returns:
//   - An error if the path is empty, content is empty, or if file I/O operations fail
//
// Example usage:
//
//	err := app.SaveFile("/path/to/file.txt", []byte("Hello, World!"))
//	if err != nil {
//	    log.Printf("Failed to save file: %v", err)
//	}
func (a *App) SaveFile(path string, content []byte) error {
	if path == "" {
		return errors.New("path is required")
	}

	if len(content) == 0 {
		return errors.New("content is required")
	}

	err := os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

// UploadFile opens a file dialog for the user to select a text file and reads its contents.
//
// This function presents a native file picker dialog that allows users to browse and
// select text files (.txt extension). The selected file is read and its contents are
// returned as a string. If no file is selected or an error occurs during file
// operations, appropriate error values are returned.
//
// Returns:
//   - A string containing the file contents if successful
//   - An empty string if no file was selected (user cancelled the dialog)
//   - An error if the file dialog fails to open or if file reading operations fail
//
// File filters:
//   - Only text files (*.txt) are shown in the file picker
//
// Example usage:
//
//	content, err := app.UploadFile()
//	if err != nil {
//	    log.Printf("Failed to upload file: %v", err)
//	} else if content == "" {
//	    log.Println("No file selected")
//	} else {
//	    log.Printf("File content: %s", content)
//	}
func (a *App) UploadFile() (string, error) {
	fd := runtime.OpenDialogOptions{
		Title: "Select a file",
		Filters: []runtime.FileFilter{
			{DisplayName: "Text files", Pattern: "*.txt"},
		},
	}
	path, err := runtime.OpenFileDialog(a.ctx, fd)
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
