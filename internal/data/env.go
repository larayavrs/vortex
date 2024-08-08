package data

import (
	"os"

	"github.com/hashicorp/go-envparse"
	"github.com/pkg/errors"
)

// The function ReadEnviromentFile reads and processes the environment configuration from a file specified by the given path.
// This function attempts to open and read the file at the provided path. The file should contain environment variables
// in a format suitable for processing (e.g., key=value pairs). The function will parse the file and apply the environment
// variables accordingly.
//
// Parameters:
//   - path: The file path to the environment configuration file.
//   - errorMissingFile: If true, the function will return an error if the file is not found. If false, the function will
//     not return an error for missing files but may handle the situation silently or with a warning.
//
// Returns:
//   - An error if the file cannot be read, if there are issues with parsing, or if the file is missing and
//     `errorMissingFile` is true. Otherwise, it returns nil indicating success.
func ReadEnviromentFile(path string, errorMissingFile bool) error {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		if errorMissingFile {
			return errors.New("Environment file not found" + path)
		}
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Failed to load environment file")
	}
	if file != nil {
		res, err := envparse.Parse(file)
		if err != nil {
			return errors.Wrap(err, "Failed to parse environment file")
		}
		for varkey, varvalue := range res {
			if _, exists := os.LookupEnv(varkey); !exists {
				if err := os.Setenv(varkey, varvalue); err != nil {
					return errors.Wrap(err, "Failed to set environment variable")
				}
			}
		}
	}
	return nil
}
