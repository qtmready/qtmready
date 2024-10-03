package defs

import (
	"encoding/json"
)

type (
	// Signal represents a signal that can be sent to a workflow. It enforces a consistent type for all workflow signals,
	// promoting code clarity and maintainability. It implements json.Marshaler and json.Unmarshaler interfaces,
	// enabling easy serialization and deserialization for communication with Temporal and other systems.
	//
	// It is recommended that each package alias this type for local constants, e.g.:
	//
	//  type Signal defs.Signal
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
