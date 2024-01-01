//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package gpu_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-hardware/ghw/pkg/gpu"
)

func TestGPU(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_GPU"); ok {
		t.Skip("Skipping GPU tests.")
	}
	if _, err := os.Stat("/sys/class/drm"); errors.Is(err, os.ErrNotExist) {
		t.Skip("Skipping GPU tests. The environment has no /sys/class/drm directory.")
	}
	info, err := gpu.New(context.TODO())
	if err != nil {
		t.Fatalf("Expected no error creating GPUInfo, but got %v", err)
	}

	if len(info.GraphicsCards) == 0 {
		t.Fatalf("Expected >0 GPU cards, but found 0.")
	}

	for _, card := range info.GraphicsCards {
		if card.Address != "" {
			di := card.DeviceInfo
			if di == nil {
				t.Fatalf("Expected card with address %s to have non-nil DeviceInfo.", card.Address)
			}
		}
		// TODO(go-hardware): Add Card.Node test when using injected sysfs for testing
	}
}
