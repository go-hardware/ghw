//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/go-hardware/ghw/pkg/baseboard"
	"github.com/go-hardware/ghw/pkg/bios"
	"github.com/go-hardware/ghw/pkg/block"
	"github.com/go-hardware/ghw/pkg/chassis"
	ghwcontext "github.com/go-hardware/ghw/pkg/context"
	"github.com/go-hardware/ghw/pkg/cpu"
	"github.com/go-hardware/ghw/pkg/gpu"
	"github.com/go-hardware/ghw/pkg/memory"
	"github.com/go-hardware/ghw/pkg/net"
	"github.com/go-hardware/ghw/pkg/pci"
	pciaddress "github.com/go-hardware/ghw/pkg/pci/address"
	"github.com/go-hardware/ghw/pkg/product"
	"github.com/go-hardware/ghw/pkg/topology"
)

type Options = ghwcontext.Options

var (
	NewContext               = ghwcontext.New
	WithRootMountpoint       = ghwcontext.WithRootMountpoint
	WithDisableWarnings      = ghwcontext.WithDisableWarnings
	WithDisableExternalTools = ghwcontext.WithDisableExternalTools
	WithOptions              = ghwcontext.WithOptions
)

type CPUInfo = cpu.Info

var (
	CPU = cpu.New
)

type MemoryArea = memory.Area
type MemoryInfo = memory.Info
type MemoryCache = memory.Cache
type MemoryCacheType = memory.CacheType
type MemoryModule = memory.Module

const (
	MemoryCacheTypeUnified     = memory.CacheTypeUnified
	MemoryCacheTypeInstruction = memory.CacheTypeInstruction
	MemoryCacheTypeData        = memory.CacheTypeData
)

var (
	Memory = memory.New
)

type BlockInfo = block.Info
type Disk = block.Disk
type Partition = block.Partition

var (
	Block = block.New
)

type DriveType = block.DriveType

const (
	DriveTypeUnknown = block.DriveTypeUnknown
	DriveTypeHDD     = block.DriveTypeHDD
	DriveTypeFDD     = block.DriveTypeFDD
	DriveTypeODD     = block.DriveTypeODD
	DriveTypeSSD     = block.DriveTypeSSD
)

type StorageController = block.StorageController

const (
	StorageControllerUnknown = block.StorageControllerUnknown
	StorageControllerIDE     = block.StorageControllerIDE
	StorageControllerSCSI    = block.StorageControllerSCSI
	StorageControllerNVMe    = block.StorageControllerNVMe
	StorageControllerVirtIO  = block.StorageControllerVirtIO
	StorageControllerMMC     = block.StorageControllerMMC
)

type NetworkInfo = net.Info
type NIC = net.NIC
type NICCapability = net.NICCapability

var (
	Network = net.New
)

type BIOSInfo = bios.Info

var (
	BIOS = bios.New
)

type ChassisInfo = chassis.Info

var (
	Chassis = chassis.New
)

type BaseboardInfo = baseboard.Info

var (
	Baseboard = baseboard.New
)

type TopologyInfo = topology.Info
type TopologyNode = topology.Node

var (
	Topology = topology.New
)

type Architecture = topology.Architecture

const (
	ArchitectureSMP  = topology.ArchitectureSMP
	ArchitectureNUMA = topology.ArchitectureNUMA
)

type PCIInfo = pci.Info
type PCIAddress = pciaddress.Address
type PCIDevice = pci.Device

var (
	PCI                  = pci.New
	PCIAddressFromString = pciaddress.FromString
)

type ProductInfo = product.Info

var (
	Product = product.New
)

type GPUInfo = gpu.Info
type GraphicsCard = gpu.GraphicsCard

var (
	GPU = gpu.New
)
