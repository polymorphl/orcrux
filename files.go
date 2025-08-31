package main

import (
	"errors"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

// SaveFileDialog opens a file dialog for the user to choose where to save a file and writes the content.
//
// This function presents a native file save dialog that allows users to browse and
// choose where to save their file. The selected path is used to write the provided
// content. If no path is selected or an error occurs during file operations,
// appropriate error values are returned.
//
// Parameters:
//   - content: The byte array to write to the file (cannot be empty)
//   - defaultName: The default filename to suggest in the save dialog
//
// Returns:
//   - An error if content is empty, no path was selected, or if file I/O operations fail
//
// File filters:
//   - Only text files (*.txt) are shown in the file picker
//
// Example usage:
//
//	err := app.SaveFileDialog([]byte("Hello, World!"), "shards.txt")
//	if err != nil {
//	    log.Printf("Failed to save file: %v", err)
//	}
func (a *App) SaveFileDialog(content []byte, defaultName string) error {
	if len(content) == 0 {
		return errors.New("content is required")
	}

	fd := runtime.SaveDialogOptions{
		Title:           "Save shards file",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{DisplayName: "Text files", Pattern: "*.txt"},
		},
	}

	path, err := runtime.SaveFileDialog(a.ctx, fd)
	if err != nil {
		return err
	}
	if path == "" {
		return nil // User cancelled the dialog
	}

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
