//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package bios

import (
	"context"
	"fmt"

	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/util"
)

// Info defines BIOS release information
type Info struct {
	// Vendor is the identifier of the BIOS vendor, if any
	Vendor string `json:"vendor"`
	// Version is the vendor-specific version of the BIOS, if any
	Version string `json:"version"`
	// Date is the date the BIOS was released
	Date string `json:"date"`
}

// String returns a human-readable description of the host's BIOS
func (i *Info) String() string {

	vendorStr := ""
	if i.Vendor != "" {
		vendorStr = " vendor=" + i.Vendor
	}
	versionStr := ""
	if i.Version != "" {
		versionStr = " version=" + i.Version
	}
	dateStr := ""
	if i.Date != "" && i.Date != util.UNKNOWN {
		dateStr = " date=" + i.Date
	}

	res := fmt.Sprintf(
		"bios%s%s%s",
		vendorStr,
		versionStr,
		dateStr,
	)
	return res
}

// New returns a pointer to a Info struct containing information
// about the host's BIOS
func New(ctx context.Context) (*Info, error) {
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate BIOS information in a top-level
// "bios" YAML/JSON map/object key
type biosPrinter struct {
	Info *Info `json:"bios"`
}

// YAMLString returns a string with the BIOS information formatted as YAML
// under a top-level "dmi:" key
func (info *Info) YAMLString() string {
	return marshal.SafeYAML(biosPrinter{info})
}

// JSONString returns a string with the BIOS information formatted as JSON
// under a top-level "bios:" key
func (info *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(biosPrinter{info}, indent)
}
