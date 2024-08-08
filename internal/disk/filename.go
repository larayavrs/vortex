package disk

import (
	"flag"
	"io"
	"os"

	"github.com/larayavrs/vortex/pkg"
	"github.com/pkg/errors"
)

// localTemplateFilenames stores the filenames of locally available templates.
// This variable holds a slice of strings where each string represents the name of a template
// file found within a specific directory or source. It is typically used to cache or reference
// the templates that have been retrieved, allowing for easier access and manipulation of these
// filenames throughout the program.
var localTemplateFilenames []string

// GetTemplateFilenames retrieves a list of template filenames from a predefined directory or source.
// This function searches for template files within a specific directory, gathers their filenames,
// and returns them as a slice of strings. The function also handles any errors that may occur
// during the process, such as issues with accessing the directory or reading the filenames.
//
// Returns:
//   - A slice of strings containing the filenames of all templates found.
//   - An error if there is an issue accessing the directory or reading the files. If no error occurs,
//     the returned error will be nil.
//
// Example usage:
//
//	filenames, err := GetTemplateFilenames()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Template filenames:", filenames)
func GetTemplateFilenames() ([]string, error) {
	if len(flag.Args()) >= 1 {
		localTemplateFilenames = flag.Args()
	}
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get file info")
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		fileNameBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read from stdin")
		}
		localTemplateFilenamesViaPipe, err := pkg.TokenizeLine(string(fileNameBytes))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to tokenize line")
		}
		if len(localTemplateFilenamesViaPipe) > 0 {
			return nil, errors.Wrap(err, "Template filenames are provided via stdin and as arguments")
		}
		localTemplateFilenames = append(localTemplateFilenames, localTemplateFilenamesViaPipe...)
	}
	return localTemplateFilenames, nil
}
