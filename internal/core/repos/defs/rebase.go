package defs

import (
	git "github.com/jeffwelling/git2go/v37"

	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	RebaseOperationKind string
	RebaseStatus        string

	RebasePayload struct {
		Rebase *eventsv1.Rebase `json:"rebase"`
		Path   string           `json:"path"`
	}

	RebaseOperation struct {
		Kind    RebaseOperationKind `json:"kind"`
		Status  RebaseStatus        `json:"status"`
		Head    string              `json:"head"`
		Message string              `json:"message"`
		Error   string              `json:"error,omitempty"`
	}

	RebaseResult struct {
		Head       string            `json:"head"`
		Status     RebaseStatus      `json:"status"`
		Conflicts  []string          `json:"conflicts"`
		Operations []RebaseOperation `json:"operations"`
		Error      string            `json:"error,omitempty"`
	}
)

const (
	RebaseStatusSuccess   RebaseStatus = "success"
	RebaseStatusFailure   RebaseStatus = "failure"
	RebaseStatusConflicts RebaseStatus = "conflicts"
	RebaseStatusUpToDate  RebaseStatus = "up_to-date"
	RebaseStatusAborted   RebaseStatus = "aborted"
	RebaseStatusPartial   RebaseStatus = "partial"
)

const (
	RebaseOperationKindPick   RebaseOperationKind = "pick"
	RebaseOperationKindReword RebaseOperationKind = "reword"
	RebaseOperationKindEdit   RebaseOperationKind = "edit"
	RebaseOperationKindSquash RebaseOperationKind = "squash"
	RebaseOperationKindFixup  RebaseOperationKind = "fixup"
)

var (
	gitOpTypeMap = map[git.RebaseOperationType]RebaseOperationKind{
		git.RebaseOperationPick:   RebaseOperationKindPick,
		git.RebaseOperationReword: RebaseOperationKindReword,
		git.RebaseOperationEdit:   RebaseOperationKindEdit,
		git.RebaseOperationSquash: RebaseOperationKindSquash,
		git.RebaseOperationFixup:  RebaseOperationKindFixup,
	}
)

func (r *RebaseResult) HasConflicts() bool {
	return len(r.Conflicts) > 0
}

func (r *RebaseResult) AddConflict(conflict string) {
	r.Conflicts = append(r.Conflicts, conflict)
}

func (r *RebaseResult) AddOperation(op git.RebaseOperationType, status RebaseStatus, head, message string, err error) {
	if err != nil {
		status = RebaseStatusFailure
	}

	err_ := err.Error()

	r.Operations = append(
		r.Operations,
		RebaseOperation{
			Kind:    gitOpTypeMap[op],
			Status:  status,
			Head:    head,
			Message: message,
			Error:   err_,
		})
}

func NewRebaseResult() *RebaseResult {
	return &RebaseResult{
		Status:     RebaseStatusFailure,
		Conflicts:  []string{},
		Operations: []RebaseOperation{},
	}
}
