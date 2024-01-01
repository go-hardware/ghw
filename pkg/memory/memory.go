//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package memory

import (
	"context"
	"fmt"
	"math"

	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/unit"
	"github.com/go-hardware/ghw/pkg/util"
)

// Module describes a single physical memory module for a host system. Pretty
// much all modern systems contain dual in-line memory modules (DIMMs).
//
// See https://en.wikipedia.org/wiki/DIMM
type Module struct {
	// Label is the system label, if any, for the memory module
	Label string `json:"label"`
	// Location stores the slot and memory channel location for the memory
	// module
	Location string `json:"location"`
	// SerialNumber is any serial number found for the memory module
	SerialNumber string `json:"serial_number"`
	// SizeBytes is the amount of physical RAM found in the memory module
	SizeBytes int64 `json:"size_bytes"`
	// Vendor contains the vendor, if any, for the memory module
	Vendor string `json:"vendor"`
}

// Area describes a set of physical memory on a host system. Non-NUMA systems
// will almost always have a single memory area containing all memory the
// system can use. NUMA systems will have multiple memory areas, one or more
// for each NUMA node/cell in the system.
type Area struct {
	// TotalPhysicalBytes is the total amount of RAM supplied by this memory
	// area
	TotalPhysicalBytes int64 `json:"total_physical_bytes"`
	// TotalUsableBytes is the total amount of RAM available for use by the
	// system from this memory area. Note that the bootloader can consume some
	// amount of memory from a memory area. The difference between
	// TotalPhysicalBytes and TotalUsableBytes is the amount of memory reserved
	// for the bootloader.
	TotalUsableBytes int64 `json:"total_usable_bytes"`
	// TotalUsedBytes is the total amount of memory consumed by the kernel and
	// all applications running on the system.
	TotalUsedBytes int64 `json:"total_used_bytes"`
	// SupportedPageSizes is a slice of sizes, in bytes, of memory pages
	// supported in this area
	SupportedPageSizes []uint64 `json:"supported_page_sizes"`
	// Modules contains a slice of `Module` pointers for any memory module
	// descriptors found for this memory area
	Modules []*Module `json:"modules"`
}

// String returns a short string with a summary of information for this memory
// area
func (a *Area) String() string {
	physs := util.UNKNOWN
	if a.TotalPhysicalBytes > 0 {
		tpb := a.TotalPhysicalBytes
		unit, unitStr := unit.AmountString(tpb)
		tpb = int64(math.Ceil(float64(a.TotalPhysicalBytes) / float64(unit)))
		physs = fmt.Sprintf("%d%s", tpb, unitStr)
	}
	usables := util.UNKNOWN
	if a.TotalUsableBytes > 0 {
		tub := a.TotalUsableBytes
		unit, unitStr := unit.AmountString(tub)
		tub = int64(math.Ceil(float64(a.TotalUsableBytes) / float64(unit)))
		usables = fmt.Sprintf("%d%s", tub, unitStr)
	}
	useds := ""
	if a.TotalUsedBytes > 0 {
		tub := a.TotalUsedBytes
		unit, unitStr := unit.AmountString(tub)
		tub = int64(math.Ceil(float64(a.TotalUsedBytes) / float64(unit)))
		useds = fmt.Sprintf(", %d%s used", tub, unitStr)
	}
	return fmt.Sprintf("memory (%s physical, %s usable%s)", physs, usables, useds)
}

// Info contains information about the memory on a host system.
type Info struct {
	Area
}

// New returns an Info struct that describes the memory on a host system.
func New(ctx context.Context) (*Info, error) {
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// String returns a short string with a summary of memory information
func (i *Info) String() string {
	return i.Area.String()
}

// simple private struct used to encapsulate memory information in a top-level
// "memory" YAML/JSON map/object key
type memoryPrinter struct {
	Info *Info `json:"memory"`
}

// YAMLString returns a string with the memory information formatted as YAML
// under a top-level "memory:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(memoryPrinter{i})
}

// JSONString returns a string with the memory information formatted as JSON
// under a top-level "memory:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(memoryPrinter{i}, indent)
}
