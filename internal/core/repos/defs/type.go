package reposdefs

import (
	"go.breu.io/durex/queues"
)

// signals.
const (
	RepoIOSignalPush queues.Signal = "repo_io__push" // signals a push event.
)
