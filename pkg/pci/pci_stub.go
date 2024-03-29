//go:build !linux
// +build !linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"context"
	"runtime"

	"github.com/pkg/errors"
)

func (i *Info) load(_ context.Context) error {
	return errors.New("pciFillInfo not implemented on " + runtime.GOOS)
}

// GetDevice returns a pointer to a Device struct that describes the PCI
// device at the requested address. If no such device could be found, returns
// nil
func (info *Info) GetDevice(_ context.Context, _ string) *Device {
	return nil
}
