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

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	pwg "github.com/sethvargo/go-password/password"

	"go.breu.io/ctrlplane/internal/shared/logger"
	"go.breu.io/ctrlplane/internal/shared/service"
	"go.breu.io/ctrlplane/internal/shared/temporal"
)

var (
	Service   *service.Service    // Global service instance. Must be initialized in main.go
	Logger    *logger.ZapAdapter  // Global logger instance. Must be initialized in main.go
	Temporal  *temporal.Temporal  // Global temporal instance. Must be initialized in main.go
	Validator *validator.Validate // Global validator instance. Must be initialized in main.go
)

// InitService initializes the global service instance.
func InitService() {
	Service = service.NewService(
		service.WithConfigFromEnv(),
		service.WithVersionFromBuildInfo(),
	)
}

// InitLogger initializes the global logger instance. The global Service instance must be initialized before calling this function.
func InitLogger(debug bool, skip int) {
	if Service == nil {
		panic("Service must be initialized before initializing logger")
	}

	Logger = logger.NewZapAdapter(logger.NewZapLogger(debug, skip), skip)
}

// InitValidator initializes the global validator instance.
func InitValidator() {
	Validator = validator.New()
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// InitServiceForTest initializes the global service instance for testing.
func InitServiceForTest() {
	Service = service.NewService(
		service.WithName("test"),
		service.WithDebug(true),
		service.WithSecret(pwg.MustGenerate(32, 8, 0, false, false)),
	)
}

func InitTemporal() {
	Temporal = temporal.NewTemporal(
		temporal.WithConfigFromEnv(),
		temporal.WithLogger(Logger),
		temporal.WithQueue(CoreQueue),
		temporal.WithQueue(ProvidersQueue),
		temporal.WithQueue(MutexQueue),
		temporal.WithClientConnection(),
	)
}
