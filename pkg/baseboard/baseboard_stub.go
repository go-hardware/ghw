//go:build !linux && !windows
// +build !linux,!windows

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package baseboard

import (
	"context"
	"runtime"

	"github.com/pkg/errors"
)

func (i *Info) load(_ context.Context) error {
	return errors.New("baseboardFillInfo not implemented on " + runtime.GOOS)
}
