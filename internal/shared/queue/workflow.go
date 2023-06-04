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

// WithWorkflowIDParent sets the parent workflow context.
func WithWorkflowIDParent(parent workflow.Context) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).parent = parent }
}

// WithWorkflowIDBlock sets the block name.
func WithWorkflowIDBlock(block string) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).block = block }
}

// WithWorkflowIDBlockID sets the block id.
func WithWorkflowIDBlockID(blockID string) WorkflowIDOption {
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

// WithWorkflowIDModifier sets the modifier name.
func WithWorkflowIDModifier(modifier string) WorkflowIDOption {
	return func(w WorkflowID) { w.(*id).mod = modifier }
}

// WithWorkflowIDModifierID sets the modifier id.
func WithWorkflowIDModifierID(modifierID string) WorkflowIDOption {
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

// String omits the empty parts and returns the formatted workflow id.
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

	sanitized := make([]string, 0)

	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			sanitized = append(sanitized, part)
		}
	}

	return strings.Join(sanitized, ".")
}

// NewWorkflowIDCreator  creates an idempotency key. Sometimes we need to signal the workflow from a completely disconnected
// part of the application. For us, it is important to arrive at the same workflow ID regardless of the conditions.
// We try to follow the block, element, modifier pattern popularized by advocates of mantainable CSS. For more info,
// https://getbem.com.
//
// Example:
// For the block github with installation id 123, the element being the repository with id 456, and the modifier being the
// pull request with id 789, we would call
//
//	id := NewWorkflowIDCreator(
//	  WithWorkflowBlock("github"),
//	  WithWorkflowBlockID("123"),
//	  WithWorkflowElement("repository"),
//	  WithWorkflowElementID("123"),
//	  WithWorkflowModifier("repository"),
//	  WithWorkflowModifierID("123"),
//	)
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
