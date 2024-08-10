package disk

import (
	"io"
	"os"
	"os/exec"

	"github.com/larayavrs/vortex/pkg"
	"github.com/pkg/errors"
)

const (
	// editFileSuffix is a constant that represents the suffix used to denote files
	// that are intended for editing. When a file is marked with this suffix, it may
	// indicate to the program that the file should be opened in an editor for modification.
	// The choice of suffix is arbitrary, but it should be unique to avoid conflicts with
	// regular filenames.
	editFileSuffix = "!"

	// fallbackEditor is a constant that specifies the default editor to use when no
	// other editor is configured or available. If a user does not have a preferred editor
	// set up, the program will fall back to this editor. In this case, Visual Studio Code
	// is set as the default fallback editor.
	fallbackEditor = "code"
)

// Function CaptureEditorOutput captures and returns the content of a temporary file after it has been edited.
// This function is used to open a temporary file in an external editor, allow the user to make edits,
// and then read the updated content of the file once the editor is closed. The function reads the entire
// content of the file and returns it as a string, along with any errors encountered during the process.
//
// Parameters:
//   - tempfile: A pointer to the `os.File` representing the temporary file that will be edited. This file
//     should already be created and accessible by the external editor.
//
// Returns:
//   - A string containing the updated content of the temporary file after editing.
//   - An error if there is an issue with opening the file, reading its content, or interacting with the editor.
func CaptureEditorOutput(tempfile *os.File) (string, error) {
	editorEnviromentVar := "VISUAL"
	editorEnvStr := os.Getenv(editorEnviromentVar)
	if editorEnvStr == "" {
		editorEnviromentVar = "EDITOR"
		editorEnvStr = os.Getenv(editorEnviromentVar)
	}
	if editorEnvStr == "" {
		_, err := exec.LookPath(fallbackEditor)
		if err != nil {
			return "", errors.New("Could not find a suitable editor to open the file. Please set the VISUAL or EDITOR environment variable to specify an editor.")
		}
		editorEnviromentVar = fallbackEditor
		editorEnvStr = fallbackEditor
	}
	editorCmdArgs, err := pkg.TokenizeLine(editorEnvStr)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to parse editor command")
	}
	editorArgs := append(editorCmdArgs[1:], tempfile.Name())
	cmd := exec.Command(editorCmdArgs[0], editorArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", errors.Wrapf(err, "Failed to run editor command %s %s", editorEnviromentVar, cmd.String())
	}
	_, err = tempfile.Seek(0, 0)
	if err != nil {
		return "", errors.Wrap(err, "Failed to seek to the beginning of the file")
	}
	tempfileContents, err := io.ReadAll(tempfile)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read the content of the file")
	}
	return string(tempfileContents), nil
}
