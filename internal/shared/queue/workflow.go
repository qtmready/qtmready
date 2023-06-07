// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package queue

import (
	"strings"

	"go.temporal.io/sdk/workflow"
)

type (
	idprops map[string]string

	wrkflopts struct {
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
func WithWorkflowParent(parent workflow.Context) WorkflowOptionProvider {
	return func(w WorkflowOptions) { w.(*wrkflopts).parent = parent }
}

// WithWorkflowBlock sets the block name.
func WithWorkflowBlock(block string) WorkflowOptionProvider {
	return func(w WorkflowOptions) {
		if w.(*wrkflopts).block != "" {
			panic(NewDuplicateIdPropError("block"))
		}

		w.(*wrkflopts).block = block
	}
}

// WithWorkflowBlockID sets the block value.
func WithWorkflowBlockID(val string) WorkflowOptionProvider {
	return func(w WorkflowOptions) {
		if w.(*wrkflopts).blockID != "" {
			panic(NewDuplicateIdPropError("block id"))
		}

		w.(*wrkflopts).blockID = val
	}
}

// WithWorkflowElement sets the element name.
func WithWorkflowElement(element string) WorkflowOptionProvider {
	return func(w WorkflowOptions) {
		if w.(*wrkflopts).elm != "" {
			panic(NewDuplicateIdPropError("element"))
		}

		w.(*wrkflopts).elm = element
	}
}

// WithWorkflowElementID sets the element value.
func WithWorkflowElementID(val string) WorkflowOptionProvider {
	return func(w WorkflowOptions) {
		if w.(*wrkflopts).elmID != "" {
			panic(NewDuplicateIdPropError("element id"))
		}

		w.(*wrkflopts).elmID = val
	}
}

// WithWorkflowMod sets the modifier name.
func WithWorkflowMod(modifier string) WorkflowOptionProvider {
	return func(w WorkflowOptions) {
		if w.(*wrkflopts).mod != "" {
			panic(NewDuplicateIdPropError("modifier"))
		}

		w.(*wrkflopts).mod = modifier
	}
}

// WithWorkflowModID sets the modifier value.
func WithWorkflowModID(val string) WorkflowOptionProvider {
	return func(w WorkflowOptions) {
		if w.(*wrkflopts).modID != "" {
			panic(NewDuplicateIdPropError("modifier id"))
		}

		w.(*wrkflopts).modID = val
	}
}

// WithWorkflowProp sets the prop given a key & value.
func WithWorkflowProp(key, val string) WorkflowOptionProvider {
	return func(w WorkflowOptions) {
		w.(*wrkflopts).propOrder = append(w.(*wrkflopts).propOrder, key)
		w.(*wrkflopts).props[key] = val
	}
}

func (w *wrkflopts) IsChild() bool {
	return w.parent != nil
}

// Suffix sanitizes the suffix and returns it.
func (w *wrkflopts) Suffix() string {
	parts := []string{w.block, w.blockID, w.elm, w.elmID, w.mod, w.modID}
	for _, key := range w.propOrder {
		parts = append(parts, key, w.props[key])
	}

	sanitized := make([]string, 0)

	// removing empty strings and trimming spaces
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			sanitized = append(sanitized, part)
		}
	}

	return strings.Join(sanitized, ".")
}

func (w *wrkflopts) ParentWorkflowID() string {
	if w.parent == nil {
		panic(ErrParentNil)
	}

	return workflow.GetInfo(w.parent).WorkflowExecution.ID
}

// NewWorkflowOptions  creates the workflow ID. Sometimes we need to signal the workflow from a completely disconnected
// part of the application. For us, it is important to arrive at the same workflow ID regardless of the conditions.
// We try to follow the block, element, modifier pattern popularized by advocates of mantainable CSS. For more info,
// https://getbem.com.
//
// Example:
// For the block github with installation id 123, the element being the repository with id 456, and the modifier being the
// pull request with id 789, we would call
//
//	opts := NewWorkflowOptions(
//	  WithWorkflowBlock("github"),
//	  WithWorkflowBlockID("123"),
//	  WithWorkflowElement("repository"),
//	  WithWorkflowElementID("123"),
//	  WithWorkflowModifier("repository"),
//	  WithWorkflowModifierID("123"),
//	)
//
//	id := opts.String()
//
// Please note, that the design is work in progress and may change.
func NewWorkflowOptions(options ...WorkflowOptionProvider) WorkflowOptions {
	w := &wrkflopts{
		props:     make(idprops),
		propOrder: make([]string, 0),
	}

	for _, option := range options {
		option(w)
	}

	return w
}
