package events

import (
	"encoding/json"
)

type (
	EventVersion string // Version is the version of the event.
)

// MarshalJSON customizes the JSON encoding for EventVersion.
func (ev EventVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ev))
}

// UnmarshalJSON customizes the JSON decoding for Version.
func (ev *EventVersion) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*ev = EventVersion(s)

	return nil
}

// String returns the string representation of the Version.
func (ev EventVersion) String() string {
	return string(ev)
}

const (
	Version_0_1_0 EventVersion = "0.1.0" // version 0.1.0.
	Version_0_1_1 EventVersion = "0.1.1" // version 0.1.1.
)

const (
	// EventVersionDefault alias for the default version. This allows for easy versioning without chaniging the code base.
	EventVersionDefault = Version_0_1_0
)
