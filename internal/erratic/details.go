package erratic

type (
	// Hints represents a map of key-value pairs providing additional information about an error.
	//
	// Example:
	//
	//  info := Hints{"field": "invalid value"}
	//  fmt.Println(info) // Output: map[string]string{"field": "invalid value"}
	Hints map[string]string
)

func NewErrorDetails(args ...string) Hints {
	odd := false

	if len(args)%2 != 0 {
		odd = true
	}

	details := make(Hints)

	for i := 0; i < len(args); i += 2 {
		details[args[i]] = args[i+1]
	}

	if odd {
		details["unknown"] = args[len(args)-1]
	}

	return details
}
