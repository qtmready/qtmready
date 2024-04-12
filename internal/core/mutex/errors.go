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

package mutex

import (
	"errors"
	"fmt"
)

var (
	ErrNilContext   = errors.New("contexts not initialized")
	ErrNoResourceID = errors.New("no resource ID provided")
)

type (
	MutexError struct {
		id   string // the id of the mutex.
		kind string // kind of error. can be "acquire lock", "release lock", or "start workflow".
	}
)

func (e *MutexError) Error() string {
	return fmt.Sprintf("%s: failed to %s.", e.id, e.kind)
}

// NewAcquireLockError creates a new acquire lock error.
func NewAcquireLockError(id string) error {
	return &MutexError{id, "acquire lock"}
}

// NewReleaseLockError creates a new release lock error.
func NewReleaseLockError(id string) error {
	return &MutexError{id, "release lock"}
}

// NewPrepareMutexError creates a new start workflow error.
func NewPrepareMutexError(id string) error {
	return &MutexError{id, "prepare mutex"}
}

func NewCleanupMutexError(id string) error {
	return &MutexError{id, "cleanup mutex"}
}
