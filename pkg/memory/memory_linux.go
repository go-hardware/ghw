// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package memory

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	ghwcontext "github.com/go-hardware/ghw/pkg/context"
	ghwpath "github.com/go-hardware/ghw/pkg/path"
	"github.com/go-hardware/ghw/pkg/unit"
	"github.com/go-hardware/ghw/pkg/util"
)

const (
	warnCannotDeterminePhysicalMemory = `
Could not determine total physical bytes of memory. This may
be due to the host being a virtual machine or container with no
/var/log/syslog file or /sys/devices/system/memory directory, or
the current user may not have necessary privileges to read the syslog.
We are falling back to setting the total physical amount of memory to
the total usable amount of memory
`
)

var (
	// System log lines will look similar to the following:
	// ... kernel: [0.000000] Memory: 24633272K/25155024K ...
	regexSyslogMemline = regexp.MustCompile(`Memory:\s+\d+K\/(\d+)K`)
	// regexMemoryBlockDirname matches a subdirectory in either
	// /sys/devices/system/memory or /sys/devices/system/node/nodeX that
	// represents information on a specific memory cell/block
	regexMemoryBlockDirname = regexp.MustCompile(`memory\d+$`)
)

func (i *Info) load(ctx context.Context) error {
	paths := ghwpath.New(ctx)
	mi := memInfo{}
	if err := mi.load(paths.ProcMeminfo); err != nil {
		return err
	}
	usable := mi.totalUsableBytes()
	if usable < 1 {
		return fmt.Errorf("Could not determine total usable bytes of memory")
	}
	i.TotalUsableBytes = usable
	used := mi.totalUsedBytes()
	if used < 1 {
		return fmt.Errorf("Could not determine total used bytes of memory")
	}
	i.TotalUsedBytes = used
	tpb := memTotalPhysicalBytes(paths)
	i.TotalPhysicalBytes = tpb
	if tpb < 1 {
		ghwcontext.Warn(ctx, warnCannotDeterminePhysicalMemory)
		i.TotalPhysicalBytes = usable
	}
	i.SupportedPageSizes, _ = memorySupportedPageSizes(paths.SysKernelMMHugepages)
	return nil
}

func AreaForNode(ctx context.Context, nodeID int) (*Area, error) {
	paths := ghwpath.New(ctx)

	var err error
	var blockSizeBytes uint64
	var totPhys int64
	var totUsable int64

	mi := memInfo{}
	if err := mi.load(paths.NodeMeminfo(nodeID)); err != nil {
		return nil, err
	}
	totUsable = mi.totalUsableBytes()
	path := filepath.Join(
		paths.SysDevicesSystemNode,
		fmt.Sprintf("node%d", nodeID),
	)

	blockSizeBytes, err = memoryBlockSizeBytes(paths.SysDevicesSystemMemory)
	if err == nil {
		totPhys, err = memoryTotalPhysicalBytesFromPath(path, blockSizeBytes)
		if err != nil {
			return nil, err
		}
	} else {
		// NOTE(jaypipes): Some platforms (e.g. ARM) will not have a
		// /sys/device/system/memory/block_size_bytes file. If this is the
		// case, we set physical bytes equal to either the physical memory
		// determined from syslog or the usable bytes
		//
		// see: https://bugzilla.redhat.com/show_bug.cgi?id=1794160
		// see: https://github.com/go-hardware/ghw/issues/336
		totPhys = memTotalPhysicalBytesFromSyslog(paths)
	}

	supportedHP, err := memorySupportedPageSizes(filepath.Join(path, "hugepages"))
	if err != nil {
		return nil, err
	}

	return &Area{
		TotalPhysicalBytes: totPhys,
		TotalUsableBytes:   totUsable,
		SupportedPageSizes: supportedHP,
	}, nil
}

func memoryBlockSizeBytes(dir string) (uint64, error) {
	// get the memory block size in byte in hexadecimal notation
	blockSize := filepath.Join(dir, "block_size_bytes")

	d, err := os.ReadFile(blockSize)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(d)), 16, 64)
}

func memTotalPhysicalBytes(paths *ghwpath.Paths) (total int64) {
	defer func() {
		// fallback to the syslog file approach in case of error
		if total < 0 {
			total = memTotalPhysicalBytesFromSyslog(paths)
		}
	}()

	// detect physical memory from /sys/devices/system/memory
	dir := paths.SysDevicesSystemMemory
	blockSizeBytes, err := memoryBlockSizeBytes(dir)
	if err != nil {
		total = -1
		return total
	}

	total, err = memoryTotalPhysicalBytesFromPath(dir, blockSizeBytes)
	if err != nil {
		total = -1
	}
	return total
}

// memoryTotalPhysicalBytesFromPath accepts a directory -- either
// /sys/devices/system/memory (for the entire system) or
// /sys/devices/system/node/nodeX (for a specific NUMA node) -- and a block
// size in bytes and iterates over the sysfs memory block subdirectories,
// accumulating blocks that are "online" to determine a total physical memory
// size in bytes
func memoryTotalPhysicalBytesFromPath(dir string, blockSizeBytes uint64) (int64, error) {
	var total int64
	files, err := os.ReadDir(dir)
	if err != nil {
		return -1, err
	}
	// There are many subdirectories of /sys/devices/system/memory or
	// /sys/devices/system/node/nodeX that are named memory{cell} where {cell}
	// is a 0-based index of the memory block. These subdirectories contain a
	// state file (e.g. /sys/devices/system/memory/memory64/state that will
	// contain the string "online" if that block is active.
	for _, file := range files {
		fname := file.Name()
		// NOTE(jaypipes): we cannot rely on file.IsDir() here because the
		// memory{cell} sysfs directories are not actual directories.
		if !regexMemoryBlockDirname.MatchString(fname) {
			continue
		}
		s, err := os.ReadFile(filepath.Join(dir, fname, "state"))
		if err != nil {
			return -1, err
		}
		// if the memory block state is 'online' we increment the total with
		// the memory block size to determine the amount of physical
		// memory available on this system.
		if strings.TrimSpace(string(s)) != "online" {
			continue
		}
		total += int64(blockSizeBytes)
	}
	return total, nil
}

func memTotalPhysicalBytesFromSyslog(paths *ghwpath.Paths) int64 {
	// In Linux, the total physical memory can be determined by looking at the
	// output of dmidecode, however dmidecode requires root privileges to run,
	// so instead we examine the system logs for startup information containing
	// total physical memory and cache the results of this.
	findPhysicalKb := func(line string) int64 {
		matches := regexSyslogMemline.FindStringSubmatch(line)
		if len(matches) == 2 {
			i, err := strconv.Atoi(matches[1])
			if err != nil {
				return -1
			}
			return int64(i * 1024)
		}
		return -1
	}

	// /var/log will contain a file called syslog and 0 or more files called
	// syslog.$NUMBER or syslog.$NUMBER.gz containing system log records. We
	// search each, stopping when we match a system log record line that
	// contains physical memory information.
	logDir := paths.VarLog
	logFiles, err := os.ReadDir(logDir)
	if err != nil {
		return -1
	}
	for _, file := range logFiles {
		if strings.HasPrefix(file.Name(), "syslog") {
			fullPath := filepath.Join(logDir, file.Name())
			unzip := strings.HasSuffix(file.Name(), ".gz")
			var r io.ReadCloser
			r, err = os.Open(fullPath)
			if err != nil {
				return -1
			}
			defer util.SafeClose(r)
			if unzip {
				r, err = gzip.NewReader(r)
				if err != nil {
					return -1
				}
			}

			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				line := scanner.Text()
				size := findPhysicalKb(line)
				if size > 0 {
					return size
				}
			}
		}
	}
	return -1
}

// memInfo is a /proc/meminfo file parsed into its key:value blocks, with all
// values ending in a "kB" suffix having their values multiplied by 1024.
type memInfo map[string]int64

// load accepts a path and loads the memInfo map by parsing the supplied
// /proc/meminfo file.
//
// In Linux, /proc/meminfo or its close relative
// /sys/devices/system/node/node*/meminfo contains a set of memory-related
// amounts, with lines looking like the following:
//
// $ cat /proc/meminfo
// MemTotal:       24677596 kB
// MemFree:        21244356 kB
// MemAvailable:   22085432 kB
// ...
// HugePages_Total:       0
// HugePages_Free:        0
// HugePages_Rsvd:        0
// HugePages_Surp:        0
// ...
//
// The /sys/devices/system/node/node*/meminfo files look like this, however:
//
// Node 0 MemTotal:       24677596 kB
// Node 0 MemFree:        21244356 kB
// Node 0 MemAvailable:   22085432 kB
// ...
// Node 0 HugePages_Total:       0
// Node 0 HugePages_Free:        0
// Node 0 HugePages_Rsvd:        0
// Node 0 HugePages_Surp:        0
// ...
//
// It's worth noting that /proc/meminfo returns exact information, not
// "theoretical" information. For instance, on the above system, I have 24GB of
// RAM but MemTotal is indicating only around 23GB. This is because MemTotal
// contains the exact amount of *usable* memory after accounting for the
// kernel's resident memory size and a few reserved bits.  Please note GHW
// cares about the subset of lines shared between system-wide and per-NUMA-node
// meminfos. For more information, see:
//
//	https://www.kernel.org/doc/Documentation/filesystems/proc.txt
func (mi memInfo) load(fp string) error {
	r, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer util.SafeClose(r)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Node") {
			// For the /sys/devices/system/node/nodeX/meminfo files, the lines
			// all start with "Node X ". We need to strip all that off.
			fields := strings.Fields(line)
			line = strings.Join(fields[2:], "")
		}
		parts := strings.Split(line, ":")
		key := strings.TrimSpace(parts[0])
		raw := parts[1]
		inKb := strings.HasSuffix(raw, "kB")
		v, err := strconv.Atoi(
			strings.TrimSpace(
				strings.TrimSuffix(
					raw, "kB",
				),
			),
		)
		v64 := int64(v)
		if err != nil {
			return err
		}
		if inKb {
			v64 = v64 * unit.KB
		}
		mi[key] = v64
	}
	return nil
}

// totalUsageBytes returns the MemTotal entry from the memInfo map
func (mi memInfo) totalUsableBytes() int64 {
	v, ok := mi["MemTotal"]
	if !ok {
		return -1
	}
	return v
}

// totalUsedBytes returns the total used memory from the memInfo map.
// We calculate used memory with the following formula:
// mem_total - (mem_free + mem_buffered + mem_cached + mem_slab_reclaimable)
func (mi memInfo) totalUsedBytes() int64 {
	mf, ok := mi["MemFree"]
	if !ok {
		return -1
	}
	mc, ok := mi["Cached"]
	if !ok {
		return -1
	}
	mb, ok := mi["Buffers"]
	if !ok {
		return -1
	}
	mt, ok := mi["MemTotal"]
	if !ok {
		return -1
	}
	if sr, ok := mi["SReclaimable"]; ok {
		return mt - (mf + mb + mc + sr)
	} else if st, ok := mi["Slab"]; ok {
		// If detailed slab information isn't present, fall back to slab total.
		return mt - (mf + mb + mc + st)
	}
	return -1
}

func memorySupportedPageSizes(hpDir string) ([]uint64, error) {
	// In Linux, /sys/kernel/mm/hugepages contains a directory per page size
	// supported by the kernel. The directory name corresponds to the pattern
	// 'hugepages-{pagesize}kb'
	out := make([]uint64, 0)

	files, err := os.ReadDir(hpDir)
	if err != nil {
		return out, err
	}
	for _, file := range files {
		parts := strings.Split(file.Name(), "-")
		sizeStr := parts[1]
		// Cut off the 'kb'
		sizeStr = sizeStr[0 : len(sizeStr)-2]
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return out, err
		}
		out = append(out, uint64(size*int(unit.KB)))
	}
	return out, nil
}
