package defs

import (
	"encoding/json"
)

type (
	// Signal is a string alias intended for defining groups of workflow signals, for example, "on_push", "on_pr", etc.
	// It ensures consistency and code clarity. The Signal type provides methods for conversion and serialization,
	// promoting good developer experience.
	//
	// NOTE: Should we rename this type to TextField or something similar? "Signal" is a bit specific, and the helper
	// methods for conversion and serialization can be used whenever we need to define a group of constants.
	Signal string
)

// String returns the string representation of a Signal.
func (w Signal) String() string {
	return string(w)
}

// MarshalJSON implements the json.Marshaler interface for Signal.
//
// It marshals the Signal as a JSON string.
func (w Signal) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(w))
}

// UnmarshalJSON implements the json.Unmarshaler interface for Signal.
//
// It unmarshals a JSON string into a Signal.
func (w *Signal) UnmarshalJSON(data []byte) error {
	var signal string
	if err := json.Unmarshal(data, &signal); err != nil {
		return err
	}

	*w = Signal(signal)

	return nil
}
