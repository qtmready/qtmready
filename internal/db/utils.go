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

package db

import (
	"strings"
	"testing"

	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid/v4"

	"go.breu.io/ctrlplane/internal/shared"
)

func CreateSlug(s string) string {
	return slug.Make(s) + "-" + strings.ToLower(shortuuid.New())
}

func TestEntityGetTable(expect string, entity Entity) func(*testing.T) {
	return func(t *testing.T) {
		if expect != entity.GetTable().Metadata().M.Name {
			t.Errorf("expected %s, got %s", expect, entity.GetTable().Metadata().M.Name)
		}
	}
}

func TestEntityOps(entity Entity, tests shared.TestFnMap) func(*testing.T) {
	return func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, test.Run(test.Args, test.Want))
		}
	}
}
