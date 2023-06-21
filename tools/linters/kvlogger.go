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

package linters

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	kvdoc   = `checks if the logger is being used with the correct key-value pairs`
	kverror = `expected a message with optional key-value pairs`
)

var (
	KVLogger = &analysis.Analyzer{
		Name:     "kvlogger",
		Doc:      kvdoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
)

func run(pass *analysis.Pass) (any, error) {
	callexprs := []ast.Node{(*ast.CallExpr)(nil)}
	scan, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	if !ok {
		return nil, fmt.Errorf("kvlogger: couldn't get inspector")
	}

	scan.Preorder(callexprs, check(pass))

	return nil, nil
}

func check(pass *analysis.Pass) func(ast.Node) {
	return func(node ast.Node) {
		function := node.(*ast.CallExpr)
		if selector, ok := function.Fun.(*ast.SelectorExpr); ok {
			// check shared.Logger().{Debug,Error,Warn,Info}
			if callexpr, ok := selector.X.(*ast.CallExpr); ok {
				if checkFnName(selector.Sel.Name) && checkIdentOrder(callexpr) && isEven(len(function.Args)) {
					pass.Reportf(function.Pos(), kverror)
				}
			}

			// check logger.{Debug,Error,Warn,Info}
			if ident, ok := selector.X.(*ast.Ident); ok {
				if ident.Name == "logger" && checkFnName(selector.Sel.Name) && isEven(len(function.Args)) {
					pass.Reportf(function.Pos(), kverror)
				}
			}
		}
	}
}

// isEven checks if the given number is even.
func isEven(n int) bool {
	return n%2 == 0
}

// checkFnName checks if the function name is Debug, Error, Warn, or Info.
func checkFnName(name string) bool {
	switch name {
	case "Debug", "Error", "Warn", "Info":
		return true
	}

	return false
}

// checkIdentOrder checks if the function is part of shared.Logger().
func checkIdentOrder(callexpr *ast.CallExpr) bool {
	selector, ok := callexpr.Fun.(*ast.SelectorExpr)
	if ok {
		ident, ok := selector.X.(*ast.Ident)
		if ok && isLogger(selector.Sel.Name) && isShared(ident.Name) {
			return true
		}
	}

	return false
}

// isLogger checks if the function name is prefixed with Logger.
func isLogger(name string) bool {
	if name == "Logger" {
		return true
	}

	return false
}

// isShared checks if the function name is prefixed with shared.
func isShared(name string) bool {
	if name == "shared" {
		return true
	}

	return false
}
