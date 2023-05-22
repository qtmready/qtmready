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
	"sync"

	"github.com/go-playground/validator/v10"
	pwg "github.com/sethvargo/go-password/password"

	"go.breu.io/ctrlplane/internal/shared/logger"
	"go.breu.io/ctrlplane/internal/shared/service"
	"go.breu.io/ctrlplane/internal/shared/temporal"
)

var (
	svc     service.Service // Global service instance.
	svcOnce sync.Once       // Global service instance initializer.

	lgr     logger.Logger // Global logger instance.
	lgrOnce sync.Once     // Global logger instance initializer

	vld     *validator.Validate // Global validator instance.
	vldOnce sync.Once           // Global validator instance initializer

	tmprl     temporal.Temporal
	tmprlOnce sync.Once
)

// Service returns the global service instance. If the global service instance has not been initialized, it will be initialized with
// default values. The benefit is, you don't need to initialize the service instance in main.go if you don't need to override the
// default values.
func Service() service.Service {
	if svc == nil {
		svcOnce.Do(func() {
			svc = service.NewService(
				service.WithConfigFromEnv(),
				service.WithVersionFromBuildInfo(),
			)
		})
	}

	return svc
}

func Logger() logger.Logger {
	if lgr == nil {
		lgrOnce.Do(func() {
			lgr = logger.NewZapAdapter(Service().GetDebug(), Service().GetLogSkipper())
		})
	}

	return lgr
}

func Validator() *validator.Validate {
	if vld == nil {
		vldOnce.Do(func() {
			vld = validator.New()
			vld.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			})
		})
	}

	return vld
}

func Temporal() temporal.Temporal {
	if tmprl == nil {
		tmprlOnce.Do(func() {
			tmprl = temporal.NewTemporal(
				temporal.WithConfigFromEnv(),
				temporal.WithLogger(Logger()),
				temporal.WithQueue(CoreQueue),
				temporal.WithQueue(ProvidersQueue),
				temporal.WithQueue(MutexQueue),
			)
		})
	}

	return tmprl
}

// InitServiceForTest initializes the global service instance for testing.
func InitServiceForTest() {
	svc = service.NewService(
		service.WithName("test"),
		service.WithDebug(true),
		service.WithSecret(pwg.MustGenerate(32, 8, 0, false, false)),
	)
}
