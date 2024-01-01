//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package gpu

import (
	"context"
	"fmt"

	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/pci"
	"github.com/go-hardware/ghw/pkg/topology"
)

// GraphicsCard describes a single graphics card on the host system
type GraphicsCard struct {
	// the PCI address where the graphics card can be found
	Address string `json:"address"`
	// The "index" of the card on the bus (generally not useful information,
	// but might as well include it)
	Index int `json:"index"`
	// pointer to a PCIDevice struct that describes the vendor and product
	// model, etc
	// TODO(go-hardware): Rename this field to PCI, instead of DeviceInfo
	DeviceInfo *pci.Device `json:"pci"`
	// Topology node that the graphics card is affined to. Will be nil if the
	// architecture is not NUMA.
	Node *topology.Node `json:"node,omitempty"`
}

// String returns a human-readable description of the graphics card
func (card *GraphicsCard) String() string {
	deviceStr := card.Address
	if card.DeviceInfo != nil {
		deviceStr = card.DeviceInfo.String()
	}
	nodeStr := ""
	if card.Node != nil {
		nodeStr = fmt.Sprintf(" [affined to NUMA node %d]", card.Node.ID)
	}
	return fmt.Sprintf(
		"card #%d %s@%s",
		card.Index,
		nodeStr,
		deviceStr,
	)
}

// Info describes the host system's GPUs/graphics cards
type Info struct {
	// GraphicsCards is a slice of pointers to `GraphicsCard` structs, one for
	// each graphics card on the host system.
	GraphicsCards []*GraphicsCard `json:"cards"`
}

// String returns a human-readable description of the host system's
// GPUs/graphics cards
func (i *Info) String() string {
	numCardsStr := "cards"
	if len(i.GraphicsCards) == 1 {
		numCardsStr = "card"
	}
	return fmt.Sprintf(
		"gpu (%d graphics %s)",
		len(i.GraphicsCards),
		numCardsStr,
	)
}

// New returns a pointer to an Info struct that contains information about the
// graphics cards on the host system
func New(ctx context.Context) (*Info, error) {
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate gpu information in a top-level
// "gpu" YAML/JSON map/object key
type gpuPrinter struct {
	Info *Info `json:"gpu"`
}

// YAMLString returns a string with the gpu information formatted as YAML
// under a top-level "gpu:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(gpuPrinter{i})
}

// JSONString returns a string with the gpu information formatted as JSON
// under a top-level "gpu:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(gpuPrinter{i}, indent)
}
