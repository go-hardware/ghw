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
	Label        string `json:"label"`
	Location     string `json:"location"`
	SerialNumber string `json:"serial_number"`
	SizeBytes    int64  `json:"size_bytes"`
	Vendor       string `json:"vendor"`
}

// Area describes a set of physical memory on a host system. Non-NUMA systems
// will almost always have a single memory area containing all memory the
// system can use. NUMA systems will have multiple memory areas, one or more
// for each NUMA node/cell in the system.
type Area struct {
	TotalPhysicalBytes int64 `json:"total_physical_bytes"`
	TotalUsableBytes   int64 `json:"total_usable_bytes"`
	TotalUsedBytes     int64 `json:"total_used_bytes"`
	// An array of sizes, in bytes, of memory pages supported in this area
	SupportedPageSizes []uint64  `json:"supported_page_sizes"`
	Modules            []*Module `json:"modules"`
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