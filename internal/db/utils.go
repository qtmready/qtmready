// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package db

import (
	"strings"

	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid/v4"
)

func CreateSlug(s string) string {
	return slug.Make(s) + "-" + strings.ToLower(shortuuid.New())
}
