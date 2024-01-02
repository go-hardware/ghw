//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package topology

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-hardware/ghw/pkg/cpu"
	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/memory"
)

// Architecture describes the overall hardware architecture. It can be either
// Symmetric Multi-Processor (SMP) or Non-Uniform Memory Access (NUMA)
type Architecture int

const (
	// ArchitectureSMP is a Symmetric Multi-Processor system
	ArchitectureSMP Architecture = iota
	// ArchitectureNUMA is a Non-Uniform Memory Access system
	ArchitectureNUMA
)

var (
	architectureString = map[Architecture]string{
		ArchitectureSMP:  "SMP",
		ArchitectureNUMA: "NUMA",
	}

	// NOTE(fromani): the keys are all lowercase and do not match
	// the keys in the opposite table `architectureString`.
	// This is done because of the choice we made in
	// Architecture:MarshalJSON.
	// We use this table only in UnmarshalJSON, so it should be OK.
	stringArchitecture = map[string]Architecture{
		"smp":  ArchitectureSMP,
		"numa": ArchitectureNUMA,
	}
)

func (a Architecture) String() string {
	return architectureString[a]
}

// NOTE(go-hardware): since serialized output is as "official" as we're going to
// get, let's lowercase the string output when serializing, in order to
// "normalize" the expected serialized output
func (a Architecture) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.ToLower(a.String()))), nil
}

func (a *Architecture) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	key := strings.ToLower(s)
	val, ok := stringArchitecture[key]
	if !ok {
		return fmt.Errorf("unknown architecture: %q", key)
	}
	*a = val
	return nil
}

// Node is an abstract construct representing a collection of processors and
// various levels of memory cache that those processors share.  In a NUMA
// architecture, there are multiple NUMA nodes, abstracted here as multiple
// Node structs. In an SMP architecture, a single Node will be available in the
// Info struct and this single struct can be used to describe the levels of
// memory caching available to the single physical processor package's physical
// processor cores
type Node struct {
	// ID is the zero-based index/identifier of the Node
	ID int `json:"id"`
	// Cores is a slice of pointers to `pkg/cpu.ProcessorCore` structs for
	// processor cores in this Node
	Cores []*cpu.ProcessorCore `json:"cores"`
	// Caches is a slice of pointers to `pkg/memory.Cache` structs for memory
	// caches on this Node
	Caches []*memory.Cache `json:"caches"`
	// Distances is a slice of integer values indicating the relative distance
	// of logical processors on the host system from this Node. The zero-based
	// index of the slice is the logical processor ID. The value of the slice
	// at that index is the relative distance of that logical processor from
	// this Node.
	Distances []int `json:"distances"`
	// Memory is a pointer to the `pkg/memory.Area` struct representing
	// physical memory affined to this Node
	Memory *memory.Area `json:"memory"`
}

// String is a short human-readable description of the Node
func (n *Node) String() string {
	return fmt.Sprintf(
		"node #%d (%d cores)",
		n.ID,
		len(n.Cores),
	)
}

// Info describes the system topology for the host hardware
type Info struct {
	// Architecture is the host system architecture
	Architecture Architecture `json:"architecture"`
	// Nodes is a slice of pointers to `Node` structs describing the
	// topological units of the host system
	Nodes []*Node `json:"nodes"`
}

// New returns a pointer to an Info struct that contains information about the
// NUMA topology on the host system
func New(ctx context.Context) (*Info, error) {
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	for _, node := range info.Nodes {
		sort.Sort(memory.SortByCacheLevelTypeFirstProcessor(node.Caches))
	}
	return info, nil
}

// String returns a short, human-readable description of the host system
// topology
func (i *Info) String() string {
	archStr := "SMP"
	if i.Architecture == ArchitectureNUMA {
		archStr = "NUMA"
	}
	res := fmt.Sprintf(
		"topology %s (%d nodes)",
		archStr,
		len(i.Nodes),
	)
	return res
}

// simple private struct used to encapsulate topology information in a
// top-level "topology" YAML/JSON map/object key
type topologyPrinter struct {
	Info *Info `json:"topology"`
}

// YAMLString returns a string with the topology information formatted as YAML
// under a top-level "topology:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(topologyPrinter{i})
}

// JSONString returns a string with the topology information formatted as JSON
// under a top-level "topology:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(topologyPrinter{i}, indent)
}
