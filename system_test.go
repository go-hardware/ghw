//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"os"
	"testing"

	ghwcontext "github.com/go-hardware/ghw/pkg/context"
)

// nolint: gocyclo
func TestSystem(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_HOST"); ok {
		t.Skip("Skipping system tests.")
	}

	ctx := ghwcontext.New(ghwcontext.WithDisableWarnings())

	system, err := System(ctx)

	if err != nil {
		t.Fatalf("Expected nil error but got %v", err)
	}
	if system == nil {
		t.Fatalf("Expected non-nil system but got nil.")
	}

	mem := system.Memory
	if mem == nil {
		t.Fatalf("Expected non-nil Memory but got nil.")
	}

	tpb := mem.TotalPhysicalBytes
	if tpb < 1 {
		t.Fatalf("Expected >0 total physical memory, but got %d", tpb)
	}

	tub := mem.TotalUsableBytes
	if tub < 1 {
		t.Fatalf("Expected >0 total usable memory, but got %d", tub)
	}

	cpu := system.CPU
	if cpu == nil {
		t.Fatalf("Expected non-nil CPU, but got nil")
	}

	cores := cpu.TotalCores
	if cores < 1 {
		t.Fatalf("Expected >0 total cores, but got %d", cores)
	}

	threads := cpu.TotalThreads
	if threads < 1 {
		t.Fatalf("Expected >0 total threads, but got %d", threads)
	}

	block := system.Block
	if block == nil {
		t.Fatalf("Expected non-nil Block but got nil.")
	}

	blockTsb := block.TotalSizeBytes
	if blockTsb < 1 {
		t.Fatalf("Expected >0 total size bytes, but got %d", blockTsb)
	}

	topology := system.Topology
	if topology == nil {
		t.Fatalf("Expected non-nil Topology but got nil.")
	}

	if len(topology.Nodes) < 1 {
		t.Fatalf("Expected >0 nodes , but got %d", len(topology.Nodes))
	}

	gpu := system.GPU
	if gpu == nil {
		t.Fatalf("Expected non-nil GPU but got nil.")
	}
}
