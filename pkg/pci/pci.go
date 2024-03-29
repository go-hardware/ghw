//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jaypipes/pcidb"

	ghwcontext "github.com/go-hardware/ghw/pkg/context"
	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/topology"
	"github.com/go-hardware/ghw/pkg/util"
)

type Device struct {
	// Address is a string with the PCI address of the device
	Address string `json:"address"`
	// Vendor is the PCI Vendor code of the device
	Vendor *pcidb.Vendor `json:"vendor"`
	// Product is the PCI Product code of the device
	Product *pcidb.Product `json:"product"`
	// Revision is any revision identifier (vendor-specific) for the device
	Revision string `json:"revision"`
	// Subsystem is the PCI subsystem code of the device
	Subsystem *pcidb.Product `json:"subsystem"`
	// Class is the PCI class of the device
	Class *pcidb.Class `json:"class"`
	// Subclass is the PCI subclass of the device
	Subclass *pcidb.Subclass `json:"subclass"`
	// ProgrammingInterface is the PCI programming interface of the device
	ProgrammingInterface *pcidb.ProgrammingInterface `json:"programming_interface"`
	// Node is a pointer to a `pkg/topology.Node` struct that the PCI device is
	// affined to. Will be nil if the architecture is not NUMA.
	Node *topology.Node `json:"node,omitempty"`
	// Driver is a string containing driver information, if any, for the device
	Driver string `json:"driver"`
}

type devIdent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type devMarshallable struct {
	Driver    string   `json:"driver"`
	Address   string   `json:"address"`
	Vendor    devIdent `json:"vendor"`
	Product   devIdent `json:"product"`
	Revision  string   `json:"revision"`
	Subsystem devIdent `json:"subsystem"`
	Class     devIdent `json:"class"`
	Subclass  devIdent `json:"subclass"`
	Interface devIdent `json:"programming_interface"`
}

// NOTE(jaypipes) Device has a custom JSON marshaller because we don't want to
// serialize the entire PCIDB information for the Vendor (which includes all of
// the vendor's products, etc). Instead, we simply serialize the ID and
// human-readable name of the vendor, product, class, etc.
func (d *Device) MarshalJSON() ([]byte, error) {
	dm := devMarshallable{
		Driver:  d.Driver,
		Address: d.Address,
		Vendor: devIdent{
			ID:   d.Vendor.ID,
			Name: d.Vendor.Name,
		},
		Product: devIdent{
			ID:   d.Product.ID,
			Name: d.Product.Name,
		},
		Revision: d.Revision,
		Subsystem: devIdent{
			ID:   d.Subsystem.ID,
			Name: d.Subsystem.Name,
		},
		Class: devIdent{
			ID:   d.Class.ID,
			Name: d.Class.Name,
		},
		Subclass: devIdent{
			ID:   d.Subclass.ID,
			Name: d.Subclass.Name,
		},
		Interface: devIdent{
			ID:   d.ProgrammingInterface.ID,
			Name: d.ProgrammingInterface.Name,
		},
	}
	return json.Marshal(dm)
}

// String contains a human-readable description of the PCI device
func (d *Device) String() string {
	vendorName := util.UNKNOWN
	if d.Vendor != nil {
		vendorName = d.Vendor.Name
	}
	productName := util.UNKNOWN
	if d.Product != nil {
		productName = d.Product.Name
	}
	className := util.UNKNOWN
	if d.Class != nil {
		className = d.Class.Name
	}
	return fmt.Sprintf(
		"%s -> driver: '%s' class: '%s' vendor: '%s' product: '%s'",
		d.Address,
		d.Driver,
		className,
		vendorName,
		productName,
	)
}

// Info contains information about PCI devices on the host system
type Info struct {
	db   *pcidb.PCIDB
	arch topology.Architecture
	// Devices is a slice of `Device` structs containing information on all PCI
	// devices on the host system
	Devices []*Device
}

// String contains a human-readable description of PCI information on the host
// system.
func (i *Info) String() string {
	return fmt.Sprintf("PCI (%d devices)", len(i.Devices))
}

// New returns a pointer to an Info struct that contains information about the
// PCI devices on the host system
func New(ctx context.Context) (*Info, error) {
	// by default we don't report NUMA information;
	// we will only if are sure we are running on NUMA architecture
	info := &Info{
		arch: topology.ArchitectureSMP,
	}

	topo, err := topology.New(ctx)
	if err == nil {
		info.arch = topo.Architecture
	} else {
		ghwcontext.Warn(ctx, "error detecting system topology: %v", err)
	}
	if err = info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// lookupDevice gets a device from cached data
func (info *Info) lookupDevice(address string) *Device {
	for _, dev := range info.Devices {
		if dev.Address == address {
			return dev
		}
	}
	return nil
}

// simple private struct used to encapsulate PCI information in a top-level
// "pci" YAML/JSON map/object key
type pciPrinter struct {
	Info *Info `json:"pci"`
}

// YAMLString returns a string with the PCI information formatted as YAML
// under a top-level "pci:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(pciPrinter{i})
}

// JSONString returns a string with the PCI information formatted as JSON
// under a top-level "pci:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(pciPrinter{i}, indent)
}
