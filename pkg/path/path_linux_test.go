//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

//go:build linux
// +build linux

package path_test

import (
	"context"
	"os"

	"testing"

	ghwcontext "github.com/go-hardware/ghw/pkg/context"
	ghwpath "github.com/go-hardware/ghw/pkg/path"
)

func TestPathRoot(t *testing.T) {
	orig, origExists := os.LookupEnv("GHW_ROOT_MOUNTPOINT")
	if origExists {
		// For tests, save the original, test an override and then at the end
		// of the test, reset to the original
		defer os.Setenv("GHW_ROOT_MOUNTPOINT", orig)
		os.Unsetenv("GHW_ROOT_MOUNTPOINT")
	} else {
		defer os.Unsetenv("GHW_ROOT_MOUNTPOINT")
	}

	paths := ghwpath.New(context.TODO())

	// No environment variable is set for GHW_ROOT_MOUNTPOINT, so pathProcCpuinfo() should
	// return the default "/proc/cpuinfo"
	path := paths.ProcCpuinfo
	if path != "/proc/cpuinfo" {
		t.Fatalf("Expected pathProcCpuInfo() to return '/proc/cpuinfo' but got %s", path)
	}

	// Now set the GHW_ROOT_MOUNTPOINT environ variable and verify that pathRoot()
	// returns that value
	os.Setenv("GHW_ROOT_MOUNTPOINT", "/host")

	paths = ghwpath.New(context.TODO())

	path = paths.ProcCpuinfo
	if path != "/host/proc/cpuinfo" {
		t.Fatalf("Expected path.ProcCpuinfo to return '/host/proc/cpuinfo' but got %s", path)
	}
}

func TestPathSpecificRoots(t *testing.T) {
	ctx := ghwcontext.New(ghwcontext.WithPathOverrides(ghwcontext.PathOverrides{
		"/proc": "/host-proc",
		"/sys":  "/host-sys",
	}))

	paths := ghwpath.New(ctx)

	path := paths.ProcCpuinfo
	expectedPath := "/host-proc/cpuinfo"
	if path != expectedPath {
		t.Fatalf("Expected path.ProcCpuInfo to return %q but got %q", expectedPath, path)
	}

	path = paths.SysBusPciDevices
	expectedPath = "/host-sys/bus/pci/devices"
	if path != expectedPath {
		t.Fatalf("Expected path.SysBusPciDevices to return %q but got %q", expectedPath, path)
	}
}

func TestPathChrootAndSpecifics(t *testing.T) {
	ctx := ghwcontext.New(
		ghwcontext.WithPathOverrides(ghwcontext.PathOverrides{
			"/proc": "/host2-proc",
			"/sys":  "/host2-sys",
		}),
		ghwcontext.WithRootMountpoint("/redirect"),
	)

	paths := ghwpath.New(ctx)

	path := paths.ProcCpuinfo
	expectedPath := "/redirect/host2-proc/cpuinfo"
	if path != expectedPath {
		t.Fatalf("Expected path.ProcCpuInfo to return %q but got %q", expectedPath, path)
	}

	path = paths.SysBusPciDevices
	expectedPath = "/redirect/host2-sys/bus/pci/devices"
	if path != expectedPath {
		t.Fatalf("Expected path.SysBusPciDevices to return %q but got %q", expectedPath, path)
	}
}
