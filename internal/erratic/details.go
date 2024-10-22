package erratic

type (
	// ErrorDetails represents a map of key-value pairs providing additional information about an error.
	//
	// Example:
	//
	//     info := ErrorDetails{"field": "invalid value"}
	//     fmt.Println(info) // Output: map[string]string{"field": "invalid value"}
	ErrorDetails map[string]string
)

func NewErrorDetails(args ...string) ErrorDetails {
	odd := false

	if len(args)%2 != 0 {
		odd = true
	}

	details := make(ErrorDetails)

	for i := 0; i < len(args); i += 2 {
		details[args[i]] = args[i+1]
	}

	if odd {
		details["unknown"] = args[len(args)-1]
	}

	return details
}
