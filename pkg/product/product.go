//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package product

import (
	"context"

	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/util"
)

// Info defines product information
type Info struct {
	// Family is a PCI product family code, if any
	Family string `json:"family"`
	// Name is the name of the product, if any
	Name string `json:"name"`
	// Vendor is the identifier of the product's vendor, if any
	Vendor string `json:"vendor"`
	// SerialNumber is the serial number assigned to the product, if any
	SerialNumber string `json:"serial_number"`
	// UUID is the UUID of the product, if any
	UUID string `json:"uuid"`
	// SKU is the stock unit identifier (SKU) of the product, if any
	SKU string `json:"sku"`
	// Version is the vendor-specific version of the product, if any
	Version string `json:"version"`
}

// String is a human-readable description of the PCI product
func (i *Info) String() string {
	familyStr := ""
	if i.Family != "" {
		familyStr = " family=" + i.Family
	}
	nameStr := ""
	if i.Name != "" {
		nameStr = " name=" + i.Name
	}
	vendorStr := ""
	if i.Vendor != "" {
		vendorStr = " vendor=" + i.Vendor
	}
	serialStr := ""
	if i.SerialNumber != "" && i.SerialNumber != util.UNKNOWN {
		serialStr = " serial=" + i.SerialNumber
	}
	uuidStr := ""
	if i.UUID != "" && i.UUID != util.UNKNOWN {
		uuidStr = " uuid=" + i.UUID
	}
	skuStr := ""
	if i.SKU != "" {
		skuStr = " sku=" + i.SKU
	}
	versionStr := ""
	if i.Version != "" {
		versionStr = " version=" + i.Version
	}

	return "product" + util.ConcatStrings(
		familyStr,
		nameStr,
		vendorStr,
		serialStr,
		uuidStr,
		skuStr,
		versionStr,
	)
}

// New returns a pointer to a Info struct containing information
// about the host's product
func New(ctx context.Context) (*Info, error) {
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate product information in a top-level
// "product" YAML/JSON map/object key
type productPrinter struct {
	Info *Info `json:"product"`
}

// YAMLString returns a string with the product information formatted as YAML
// under a top-level "dmi:" key
func (info *Info) YAMLString() string {
	return marshal.SafeYAML(productPrinter{info})
}

// JSONString returns a string with the product information formatted as JSON
// under a top-level "product:" key
func (info *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(productPrinter{info}, indent)
}
