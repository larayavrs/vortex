package disk

// Defines the order of preference for backend tools used to execute HTTP requests.
// This slice of strings specifies the priority in which different backend tools should be selected
// when performing HTTP operations. The tools are listed in order of their priority, with the first
// tool in the list being the most preferred and the last one being the least preferred.
var backendPriorityOrder = []string{
	"curl",
	"httpie",
	"wget",
}

// Defines the template used to generate a starter configuration file for the program.
// This template is used to create a new configuration file with default values and placeholders.
// The template includes sections for the host, headers, query parameters, request body, and backend options.
// The user can fill in the specific values for each section to customize the configuration file.
var starterTemplate = `[Host]
http://localhost:${PORT}

[Headers]
Content-Type: application/json

# [Query]
# key1=value1&key2=value2

# [Body]
# {
#   "key": "value"
# }

[Backend]
{{ Backends }}
`
