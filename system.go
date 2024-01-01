//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"context"
	"fmt"

	"github.com/go-hardware/ghw/pkg/baseboard"
	"github.com/go-hardware/ghw/pkg/bios"
	"github.com/go-hardware/ghw/pkg/block"
	"github.com/go-hardware/ghw/pkg/chassis"
	"github.com/go-hardware/ghw/pkg/cpu"
	"github.com/go-hardware/ghw/pkg/gpu"
	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/memory"
	"github.com/go-hardware/ghw/pkg/net"
	"github.com/go-hardware/ghw/pkg/pci"
	"github.com/go-hardware/ghw/pkg/product"
	"github.com/go-hardware/ghw/pkg/topology"
)

// SystemInfo is a wrapper struct containing information about the host
// system's memory, block storage, CPU, etc
type SystemInfo struct {
	Memory    *memory.Info    `json:"memory"`
	Block     *block.Info     `json:"block"`
	CPU       *cpu.Info       `json:"cpu"`
	Topology  *topology.Info  `json:"topology"`
	Network   *net.Info       `json:"network"`
	GPU       *gpu.Info       `json:"gpu"`
	Chassis   *chassis.Info   `json:"chassis"`
	BIOS      *bios.Info      `json:"bios"`
	Baseboard *baseboard.Info `json:"baseboard"`
	Product   *product.Info   `json:"product"`
	PCI       *pci.Info       `json:"pci"`
}

// System returns a pointer to a SystemInfo struct that contains fields with
// information about the host system's CPU, memory, network devices, etc
func System(ctx context.Context) (*SystemInfo, error) {
	memInfo, err := memory.New(ctx)
	if err != nil {
		return nil, err
	}
	blockInfo, err := block.New(ctx)
	if err != nil {
		return nil, err
	}
	cpuInfo, err := cpu.New(ctx)
	if err != nil {
		return nil, err
	}
	topologyInfo, err := topology.New(ctx)
	if err != nil {
		return nil, err
	}
	netInfo, err := net.New(ctx)
	if err != nil {
		return nil, err
	}
	gpuInfo, err := gpu.New(ctx)
	if err != nil {
		return nil, err
	}
	chassisInfo, err := chassis.New(ctx)
	if err != nil {
		return nil, err
	}
	biosInfo, err := bios.New(ctx)
	if err != nil {
		return nil, err
	}
	baseboardInfo, err := baseboard.New(ctx)
	if err != nil {
		return nil, err
	}
	productInfo, err := product.New(ctx)
	if err != nil {
		return nil, err
	}
	pciInfo, err := pci.New(ctx)
	if err != nil {
		return nil, err
	}
	return &SystemInfo{
		CPU:       cpuInfo,
		Memory:    memInfo,
		Block:     blockInfo,
		Topology:  topologyInfo,
		Network:   netInfo,
		GPU:       gpuInfo,
		Chassis:   chassisInfo,
		BIOS:      biosInfo,
		Baseboard: baseboardInfo,
		Product:   productInfo,
		PCI:       pciInfo,
	}, nil
}

// String returns a newline-separated output of the SystemInfo's component
// structs' String-ified output
func (info *SystemInfo) String() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		info.Block.String(),
		info.CPU.String(),
		info.GPU.String(),
		info.Memory.String(),
		info.Network.String(),
		info.Topology.String(),
		info.Chassis.String(),
		info.BIOS.String(),
		info.Baseboard.String(),
		info.Product.String(),
		info.PCI.String(),
	)
}

// YAMLString returns a string with the host information formatted as YAML
// under a top-level "host:" key
func (i *SystemInfo) YAMLString() string {
	return marshal.SafeYAML(i)
}

// JSONString returns a string with the host information formatted as JSON
// under a top-level "host:" key
func (i *SystemInfo) JSONString(indent bool) string {
	return marshal.SafeJSON(i, indent)
}
