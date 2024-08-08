package data

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// CreateBodyTempfile creates a temporary file to store the request body.
// This method generates a temporary file with a unique name and writes
// the contents of the Body field from the RequestConfig to this file.
// The file is intended to be used for temporary storage during the request
// and may be deleted or handled according to the Tempfile field in
// the RequestConfig.
//
// Returns an error if the file creation or writing process fails.s
func (rc *RequestConfig) CreateBodyTempfile() error {
	// Check if the Body field is empty
	if len(rc.Body) == 0 {
		return nil
	}
	tmpfile_dir := ""
	if rc.Verbose {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "Failed to get current working directory")
		}
		// Assume that the temporary file will be created in the current working directory
		tmpfile_dir = cwd
	}
	// Create a temporary file with a unique name
	bodystr := strings.Join(rc.Body, "\n")
	tmpfile, err := os.CreateTemp(tmpfile_dir, "vortex-body")
	if err != nil {
		return errors.Wrap(err, "Failed to create temporary file")
	}
	if _, err := tmpfile.Write([]byte(bodystr)); err != nil {
		_ = tmpfile.Close()
		return errors.Wrap(err, "Failed to write to temporary file")
	}
	rc.TempfileName = tmpfile.Name()
	// Close the file to ensure that it is flushed and can be read by other processes
	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "Failed to close temporary file")
	}
	return nil
}

// The function RemoveBodyTempfile removes the temporary file used to store the request body.
// This method deletes the temporary file if it exists. The behavior of the
// deletion process can be influenced by the `force` parameter and the `Tempfile`
// field in the RequestConfig. If `Tempfile` is true, the file may not be removed
// even if `force` is set to true.
//
// Parameters:
//   - force: If true, the method will attempt to remove the temporary file even if
//     `Tempfile` is set to true. If false, the file will only be removed if
//     `Tempfile` is false.
//
// Returns an error if the file removal fails, such as when the file does not exist
// or there are permission issues.
func (rc *RequestConfig) RemoveBodyTempfile(force bool) error {
	if rc.TempfileName == "" {
		return nil
	}
	if !force && rc.Tempfile {
		return nil
	}
	err := os.Remove(rc.TempfileName)
	rc.TempfileName = ""
	return errors.Wrap(err, "Failed to remove temporary file")
}
