// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the 
// Breu Community License Agreement ("BCL Agreement"), version 1.0, found at  
// https://www.breu.io/license/community. By installating, downloading, 
// accessing, using or distrubting any of the software, you agree to the  
// terms of the license agreement. 
//
// The above copyright notice and the subsequent license agreement shall be 
// included in all copies or substantial portions of the software. 
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, 
// IMPLIED, STATUTORY, OR OTHERWISE, AND SPECIFICALLY DISCLAIMS ANY WARRANTY OF 
// MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE 
// SOFTWARE. 
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT 
// LIMITED TO, LOST PROFITS OR ANY CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, 
// OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, ARISING 
// OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  
// APPLICABLE LAW. 

package github

import (
	"errors"
)

var (
	ErrNoEventToParse               = errors.New("no event specified to parse")
	ErrInvalidHTTPMethod            = errors.New("invalid HTTP Method")
	ErrMissingHeaderGithubEvent     = errors.New("missing X-GitHub-Event Header")
	ErrMissingHeaderGithubSignature = errors.New("missing X-Hub-Signature Header")
	ErrInvalidEvent                 = errors.New("event not defined to be parsed")
	ErrPayloadParser                = errors.New("error parsing payload")
	ErrVerifySignature              = errors.New("HMAC verification failed")
)
