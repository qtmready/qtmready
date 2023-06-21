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

package shared

type (
	logger interface {
		Debug(string, ...any)
		Error(string, ...any)
		Info(string, ...any)
		Printf(string, ...any)
		Sync() error
		Trace(string, ...any)
		Warn(string, ...any)
		Verbose() bool
	}

	log struct{}
)

func (l *log) Debug(msg string, fields ...any)  {}
func (l *log) Error(msg string, fields ...any)  {}
func (l *log) Info(msg string, fields ...any)   {}
func (l *log) Printf(msg string, fields ...any) {}
func (l *log) Sync() error                      { return nil }
func (l *log) Trace(msg string, fields ...any)  {}
func (l *log) Warn(msg string, fields ...any)   {}
func (l *log) Verbose() bool                    { return false }

func Logger() logger {
	return &log{}
}
