//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package chassis

import (
	"context"

	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/util"
)

var (
	chassisTypeDescriptions = map[string]string{
		"1":  "Other",
		"2":  "Unknown",
		"3":  "Desktop",
		"4":  "Low profile desktop",
		"5":  "Pizza box",
		"6":  "Mini tower",
		"7":  "Tower",
		"8":  "Portable",
		"9":  "Laptop",
		"10": "Notebook",
		"11": "Hand held",
		"12": "Docking station",
		"13": "All in one",
		"14": "Sub notebook",
		"15": "Space-saving",
		"16": "Lunch box",
		"17": "Main server chassis",
		"18": "Expansion chassis",
		"19": "SubChassis",
		"20": "Bus Expansion chassis",
		"21": "Peripheral chassis",
		"22": "RAID chassis",
		"23": "Rack mount chassis",
		"24": "Sealed-case PC",
		"25": "Multi-system chassis",
		"26": "Compact PCI",
		"27": "Advanced TCA",
		"28": "Blade",
		"29": "Blade enclosure",
		"30": "Tablet",
		"31": "Convertible",
		"32": "Detachable",
		"33": "IoT gateway",
		"34": "Embedded PC",
		"35": "Mini PC",
		"36": "Stick PC",
	}
)

// Info defines chassis release information
type Info struct {
	// AssetTag is the asset tag assigned to the chassis, if any
	AssetTag string `json:"asset_tag"`
	// SerialNumber is the serial number assigned to the chassis, if any
	SerialNumber string `json:"serial_number"`
	// Type is the type code of the chassis
	Type string `json:"type"`
	// TypeDescription is the long-form string of the type code of the chassis
	TypeDescription string `json:"type_description"`
	// Vendor is the identifier of the chassis's vendor, if any
	Vendor string `json:"vendor"`
	// Version is the vendor-specific version of the chassis, if any
	Version string `json:"version"`
}

func (i *Info) String() string {
	vendorStr := ""
	if i.Vendor != "" {
		vendorStr = " vendor=" + i.Vendor
	}
	serialStr := ""
	if i.SerialNumber != "" && i.SerialNumber != util.UNKNOWN {
		serialStr = " serial=" + i.SerialNumber
	}
	versionStr := ""
	if i.Version != "" {
		versionStr = " version=" + i.Version
	}

	return "chassis type=" + util.ConcatStrings(
		i.TypeDescription,
		vendorStr,
		serialStr,
		versionStr,
	)
}

// New returns a pointer to a Info struct containing information
// about the host's chassis
func New(ctx context.Context) (*Info, error) {
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate chassis information in a top-level
// "chassis" YAML/JSON map/object key
type chassisPrinter struct {
	Info *Info `json:"chassis"`
}

// YAMLString returns a string with the chassis information formatted as YAML
// under a top-level "dmi:" key
func (info *Info) YAMLString() string {
	return marshal.SafeYAML(chassisPrinter{info})
}

// JSONString returns a string with the chassis information formatted as JSON
// under a top-level "chassis:" key
func (info *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(chassisPrinter{info}, indent)
}
