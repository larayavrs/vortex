// Package data provides functionalities for configuring and handling data-related operations,
// including setting timeouts for requests and defining custom query delimiters.

package data

import (
	"net/url"
)

// The constants `UnsetTimeout` represents the value used to indicate that no timeout is set for a request.
// When the Timeout field in the Config struct is set to UnsetTimeout (-1),
// it means that the request should not have a timeout. This is useful when
// the request is expected to take a long time to complete.
const UnsetTimeout = -1

// Config holds the configuration settings for the data package.
// This struct is used to configure various aspects of how data is processed and retrieved.
type Config struct {
	// Timeout specifies the maximum duration, in seconds, for a request to be completed.
	// If set to UnsetTimeout (-1), the request will have no timeout and could potentially run indefinitely.
	Timeout int32

	// QueryDelim is a pointer to a string that specifies the delimiter used to separate
	// multiple query parameters in a request. If nil, a default delimiter (e.g., "&") may be used.
	QueryDelim *string
}

// The function `NewConfig` creates and returns a new `Config` instance with default settings.
// By default, the `Timeout` field is set to `UnsetTimeout`, indicating that no timeout is configured.
// This function provides a convenient way to initialize a `Config` struct with default values,
// which can then be customized as needed.
func NewConfig() Config {
	return Config{
		Timeout: UnsetTimeout,
	}
}

// RequestConfig contains the necessary configuration and data for making a request to a backend service.
// This struct holds information such as the target URL, HTTP method, headers, and other relevant options.
type RequestConfig struct {
	// Host is the URL of the backend service to which the request will be sent.
	// This field is required and must be a valid URL.
	Host *url.URL

	// Body contains the lines of text that will be sent as the body of the request.
	// This can be used for sending data in a POST, PUT, or similar HTTP request.
	Body []string

	// Method specifies the HTTP method to be used for the request, such as "GET", "POST", "PUT", etc.
	// It determines the action to be performed on the resource identified by the Host.
	Method string

	// Headers contains the HTTP headers that will be included with the request.
	// These headers can be used to provide additional information such as content type or authorization tokens.
	Headers []string

	// Backend specifies the name or type of the backend service being used.
	// This could refer to a specific service or API that is being called.
	Backend string

	// BackendOptions holds additional options for configuring the backend service.
	// These options are represented as a slice of slices, where each inner slice contains related options.
	BackendOptions [][]string

	// Verbose, if true, will output the command used to perform the request.
	// This can be useful for debugging or logging the exact request being made.
	Verbose bool

	// Tempfile, if true, prevents the deletion of any temporary files generated during the request.
	// This can be useful if the temporary file needs to be inspected or reused.
	Tempfile bool

	// TempfileName specifies the name of the temporary file that will be used during the request.
	// If a temporary file is required, this name will be used, and the file will be created and managed accordingly.
	TempfileName string
}

// The type TimeoutContextValueKey is an empty struct used as a key for storing and retrieving
// timeout-related values from a context.Context. It serves as a unique identifier
// for the timeout value to avoid conflicts with other context values.
//
// This type is typically used in conjunction with the context package to store and
// access timeout information in a type-safe manner. Since the struct is empty,
// it only serves as a unique key and does not carry any data itself.
type TimeoutContextValueKey struct{}

// The type RequestResult contains the result of executing a request or command, including standard output,
// standard error, and the exit code. This struct is used to capture and process the results
// of backend operations or commands.
type RequestResult struct {
	// Stderr captures the standard error output produced during the execution.
	Stderr string

	// Stdout captures the standard output produced during the execution.
	Stdout string

	// ExitCode represents the exit status code of the executed command or process.
	// A value of 0 typically indicates success, while non-zero values indicate errors.
	ExitCode int
}
