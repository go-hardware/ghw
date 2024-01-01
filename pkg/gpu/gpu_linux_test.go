//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package gpu_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-hardware/ghw/cmd/ghw/snapshot"
	ghwcontext "github.com/go-hardware/ghw/pkg/context"
	"github.com/go-hardware/ghw/pkg/gpu"

	"github.com/go-hardware/ghw/testdata"
)

// testcase for https://github.com/go-hardware/ghw/issues/234 if nothing else:
// demonstrate how to consume snapshots from tests; test a boundary condition
// actually happened in the wild, even though on a VM environment.
func TestGPUWithoutNUMANodeInfo(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_GPU"); ok {
		t.Skip("Skipping GPU tests.")
	}

	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	workstationSnapshot := filepath.Join(
		testdataPath,
		"linux-amd64-amd-ryzen-1600.tar.gz",
	)
	// from now on we use constants reflecting the content of the snapshot we
	// requested, which we reviewed beforehand. IOW, you need to know the
	// content of the snapshot to fully understand this test. Inspect it using
	// ghw gpu -s "/path/to/linux-amd64-amd-ryzen-1600.tar.gz"

	toPath := t.TempDir()
	err = snapshot.Expand(workstationSnapshot, toPath)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	err = os.Remove(
		filepath.Join(toPath, "/sys/class/drm/card0/device/numa_node"),
	)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Cannot remove the NUMA node info: %v", err)
	}

	ctx := ghwcontext.New(
		ghwcontext.WithDisableWarnings(),
		ghwcontext.WithRootMountpoint(toPath),
	)
	info, err := gpu.New(ctx)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil GPUInfo, but got nil")
	}
	if len(info.GraphicsCards) == 0 {
		t.Fatalf("Expected >0 GPU cards, but found 0.")
	}
}
