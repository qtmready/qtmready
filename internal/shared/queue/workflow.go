package queue

import (
	"fmt"
	"strings"

	"go.temporal.io/sdk/workflow"
)

type (
	idprops map[string]string

	id struct {
		parent    workflow.Context
		block     string
		blockID   string
		elm       string
		elmID     string
		mod       string
		modID     string
		props     idprops
		propOrder []string
	}
)

// WithWorkflowParent sets the parent workflow context.
func WithWorkflowParent(parent workflow.Context) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).parent = parent }
}

// WithWorkflowBlock sets the block name.
func WithWorkflowBlock(block string) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).block = block }
}

// WithWorkflowBlockID sets the block id.
func WithWorkflowBlockID(blockID string) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).blockID = blockID }
}

// WithWorkflowElement sets the element name.
func WithWorkflowElement(element string) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).elm = element }
}

// WithWorkflowElementID sets the element id.
func WithWorkflowElementID(elementID string) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).elmID = elementID }
}

// WithWorkflowModifier sets the modifier name.
func WithWorkflowModifier(modifier string) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).mod = modifier }
}

// WithWorkflowModifierID sets the modifier id.
func WithWorkflowModifierID(modifierID string) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).modID = modifierID }
}

// WithWorkflowIDProp sets the workflow id properties.
func WithWorkflowIDProp(key, val string) WorkflowIDOption {
	return func(w WorkflowID) {
		w.(*id).propOrder = append(w.(*id).propOrder, key)
		w.(*id).props[key] = val
	}
}

func (w *id) IsChild() bool {
	return w.parent != nil
}

func (w *id) String(queue Queue) string {
	if w.parent == nil && queue == nil {
		panic("parent workflow context and queue cannot both be nil")
	}

	parent := ""
	if w.parent != nil {
		parent = workflow.GetInfo(w.parent).WorkflowExecution.ID
	} else {
		parent = fmt.Sprintf("%s.%s", queue.Prefix(), queue.Name())
	}

	parts := []string{parent, w.block, w.blockID, w.elm, w.elmID, w.mod, w.modID}
	for _, key := range w.propOrder {
		parts = append(parts, key, w.props[key])
	}

	cleaned := make([]string, 0)

	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			cleaned = append(cleaned, part)
		}
	}

	return strings.Join(parts, ".")
}

// NewWorkflowIDCreator creates a new WorkflowID.
func NewWorkflowIDCreator(options ...WorkflowIDOption) WorkflowID {
	w := &id{
		props:     make(idprops),
		propOrder: make([]string, 0),
	}

	for _, option := range options {
		option(w)
	}

	return w
}

func format(args ...string) string {
	return strings.Join(args, ".")
}
