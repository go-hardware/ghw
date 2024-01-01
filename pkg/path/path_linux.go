// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package path

import (
	"context"
	"fmt"
	"path/filepath"

	ghwcontext "github.com/go-hardware/ghw/pkg/context"
)

// PathRoots holds the roots of all the filesystem subtrees
// ghw wants to access.
type PathRoots struct {
	Etc  string
	Proc string
	Run  string
	Sys  string
	Var  string
}

// DefaultPathRoots return the canonical default value for PathRoots
func DefaultPathRoots() PathRoots {
	return PathRoots{
		Etc:  "/etc",
		Proc: "/proc",
		Run:  "/run",
		Sys:  "/sys",
		Var:  "/var",
	}
}

// PathRootsFromContext initialize PathRoots from the given Context,
// allowing overrides of the canonical default paths.
func PathRootsFromContext(ctx context.Context) PathRoots {
	roots := DefaultPathRoots()
	opts := ghwcontext.OptionsFromContext(ctx)
	if opts.PathOverrides == nil {
		return roots
	}
	overrides := *opts.PathOverrides
	if p, ok := overrides["/etc"]; ok {
		roots.Etc = p
	}
	if p, ok := overrides["/proc"]; ok {
		roots.Proc = p
	}
	if p, ok := overrides["/run"]; ok {
		roots.Run = p
	}
	if p, ok := overrides["/sys"]; ok {
		roots.Sys = p
	}
	if p, ok := overrides["/var"]; ok {
		roots.Var = p
	}
	return roots
}

type Paths struct {
	VarLog                 string
	ProcMeminfo            string
	ProcCpuinfo            string
	ProcMounts             string
	SysKernelMMHugepages   string
	SysBlock               string
	SysDevicesSystemNode   string
	SysDevicesSystemMemory string
	SysDevicesSystemCPU    string
	SysBusPciDevices       string
	SysClassDRM            string
	SysClassDMI            string
	SysClassNet            string
	RunUdevData            string
}

// New returns a new Paths struct containing filepath fields relative to the
// supplied Context
func New(ctx context.Context) *Paths {
	opts := ghwcontext.OptionsFromContext(ctx)
	root := *opts.RootMountpoint
	roots := PathRootsFromContext(ctx)
	return &Paths{
		VarLog:                 filepath.Join(root, roots.Var, "log"),
		ProcMeminfo:            filepath.Join(root, roots.Proc, "meminfo"),
		ProcCpuinfo:            filepath.Join(root, roots.Proc, "cpuinfo"),
		ProcMounts:             filepath.Join(root, roots.Proc, "self", "mounts"),
		SysKernelMMHugepages:   filepath.Join(root, roots.Sys, "kernel", "mm", "hugepages"),
		SysBlock:               filepath.Join(root, roots.Sys, "block"),
		SysDevicesSystemNode:   filepath.Join(root, roots.Sys, "devices", "system", "node"),
		SysDevicesSystemMemory: filepath.Join(root, roots.Sys, "devices", "system", "memory"),
		SysDevicesSystemCPU:    filepath.Join(root, roots.Sys, "devices", "system", "cpu"),
		SysBusPciDevices:       filepath.Join(root, roots.Sys, "bus", "pci", "devices"),
		SysClassDRM:            filepath.Join(root, roots.Sys, "class", "drm"),
		SysClassDMI:            filepath.Join(root, roots.Sys, "class", "dmi"),
		SysClassNet:            filepath.Join(root, roots.Sys, "class", "net"),
		RunUdevData:            filepath.Join(root, roots.Run, "udev", "data"),
	}
}

func (p *Paths) NodeMeminfo(nodeID int) string {
	return filepath.Join(
		p.SysDevicesSystemNode,
		fmt.Sprintf("node%d", nodeID),
		"meminfo",
	)
}

func (p *Paths) NodeCPU(nodeID int, lpID int) string {
	return filepath.Join(
		p.SysDevicesSystemNode,
		fmt.Sprintf("node%d", nodeID),
		fmt.Sprintf("cpu%d", lpID),
	)
}

func (p *Paths) NodeCPUCache(nodeID int, lpID int) string {
	return filepath.Join(
		p.NodeCPU(nodeID, lpID),
		"cache",
	)
}

func (p *Paths) NodeCPUCacheIndex(nodeID int, lpID int, cacheIndex int) string {
	return filepath.Join(
		p.NodeCPUCache(nodeID, lpID),
		fmt.Sprintf("index%d", cacheIndex),
	)
}
