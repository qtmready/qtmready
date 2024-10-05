package defs

import (
	"go.breu.io/durex/queues"
)

type (
	// Signal is a string alias intended for defining groups of workflow signals, for example, "on_push", "on_pr", etc.
	// It ensures consistency and code clarity. The Signal type provides methods for conversion and serialization,
	// promoting good developer experience.
	//
	// NOTE: Should we rename this type to TextField or something similar? "Signal" is a bit specific, and the helper
	// methods for conversion and serialization can be used whenever we need to define a group of constants.
	Signal = queues.WorkflowSignal
)
