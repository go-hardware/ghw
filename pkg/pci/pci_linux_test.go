//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-hardware/ghw/cmd/ghw/snapshot"
	ghwcontext "github.com/go-hardware/ghw/pkg/context"
	"github.com/go-hardware/ghw/pkg/marshal"
	"github.com/go-hardware/ghw/pkg/pci"
	"github.com/go-hardware/ghw/pkg/util"

	"github.com/go-hardware/ghw/testdata"
)

type pciTestCase struct {
	addr     string
	node     int
	revision string
	driver   string
}

// nolint: gocyclo
func TestPCINUMANode(t *testing.T) {
	ctx, info := pciTestSetup(t)

	tCases := []pciTestCase{
		{
			addr: "0000:07:03.0",
			// -1 is actually what we get out of the box on the snapshotted box
			node: -1,
		},
		{
			addr: "0000:05:11.0",
			node: 0,
		},
		{
			addr: "0000:05:00.1",
			node: 1,
		},
	}
	for _, tCase := range tCases {
		t.Run(fmt.Sprintf("%s (%d)", tCase.addr, tCase.node), func(t *testing.T) {
			dev := info.GetDevice(ctx, tCase.addr)
			if dev == nil {
				t.Fatalf("got nil device for address %q", tCase.addr)
			}
			if dev.Node == nil {
				if tCase.node != -1 {
					t.Fatalf("got nil numa NODE for address %q", tCase.addr)
				}
			} else {
				if dev.Node.ID != tCase.node {
					t.Errorf("got NUMA node info %#v, expected on node %d", dev.Node, tCase.node)
				}
			}
		})
	}
}

// nolint: gocyclo
func TestPCIDeviceRevision(t *testing.T) {
	ctx, info := pciTestSetup(t)

	var tCases []pciTestCase = []pciTestCase{
		{
			addr:     "0000:07:03.0",
			revision: "0x0a",
		},
		{
			addr:     "0000:05:00.0",
			revision: "0x01",
		},
	}
	for _, tCase := range tCases {
		t.Run(tCase.addr, func(t *testing.T) {
			dev := info.GetDevice(ctx, tCase.addr)
			if dev == nil {
				t.Fatalf("got nil device for address %q", tCase.addr)
			}
			if dev.Revision != tCase.revision {
				t.Errorf("device %q got revision %q expected %q", tCase.addr, dev.Revision, tCase.revision)
			}
		})
	}
}

// nolint: gocyclo
func TestPCIDriver(t *testing.T) {
	ctx, info := pciTestSetup(t)

	tCases := []pciTestCase{
		{
			addr:   "0000:07:03.0",
			driver: "mgag200",
		},
		{
			addr:   "0000:05:11.0",
			driver: "igbvf",
		},
		{
			addr:   "0000:05:00.1",
			driver: "igb",
		},
	}
	for _, tCase := range tCases {
		t.Run(fmt.Sprintf("%s (%s)", tCase.addr, tCase.driver), func(t *testing.T) {
			dev := info.GetDevice(ctx, tCase.addr)
			if dev == nil {
				t.Fatalf("got nil device for address %q", tCase.addr)
			}
			if dev.Driver != tCase.driver {
				t.Errorf("got driver %q expected %q", dev.Driver, tCase.driver)
			}
		})
	}
}

func TestPCIMarshalJSON(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := pci.New(context.TODO())
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	dev := info.ParseDevice("0000:3c:00.0", "pci:v0000144Dd0000A804sv0000144Dsd0000A801bc01sc08i02\n")
	if dev == nil {
		t.Fatalf("Failed to parse valid modalias")
	}
	s := marshal.SafeJSON(dev, true)
	if s == "" {
		t.Fatalf("Error marshalling device: %v", dev)
	}
}

// the sriov-device-plugin code has a test like this
func TestPCIMalformedModalias(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := pci.New(context.TODO())
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	var dev *pci.Device
	dev = info.ParseDevice("0000:00:01.0", "pci:junk")
	if dev != nil {
		t.Fatalf("Parsed successfully junk data")
	}

	dev = info.ParseDevice("0000:00:01.0", "pci:v00008086d00005916sv000017AAsd0000224Bbc03sc00i00extrajunkextradataextraextra")
	if dev == nil {
		t.Fatalf("Failed to parse valid modalias with extra data")
	}
}

func pciTestSetup(t *testing.T) (context.Context, *pci.Info) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}

	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	multiNUMASnapshot := filepath.Join(testdataPath, "linux-amd64-intel-xeon-L5640.tar.gz")

	// from now on we use constants reflecting the content of the snapshot we
	// requested, which we reviewed beforehand. IOW, you need to know the
	// content of the snapshot to fully understand this test. Inspect it using
	// ghwc topology -s "/path/to/linux-amd64-intel-xeon-L5640.tar.gz"

	toPath := t.TempDir()
	err = snapshot.Expand(multiNUMASnapshot, toPath)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	ctx := ghwcontext.New(
		ghwcontext.WithRootMountpoint(toPath),
	)
	info, err := pci.New(ctx)

	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil PCIInfo, but got nil")
	}
	return ctx, info
}

// we have this test in pci_linux_test.go (and not in pci_test.go) because `pciFillInfo` is implemented
// only on linux; so having it in the platform-independent tests would lead to false negatives.
func TestPCIMarshalUnmarshal(t *testing.T) {
	ctx := ghwcontext.New(ghwcontext.WithDisableWarnings())
	data, err := pci.New(ctx)
	if err != nil {
		t.Fatalf("Expected no error creating pci.Info, but got %v", err)
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Expected no error marshaling pci.Info, but got %v", err)
	}

	var topo *pci.Info

	err = json.Unmarshal(jdata, &topo)
	if err != nil {
		t.Fatalf("Expected no error unmarshaling pci.Info, but got %v", err)
	}
}

func TestPCIModaliasWithUpperCaseClassID(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := pci.New(context.TODO())
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	dev := info.ParseDevice("0000:00:1f.4", "pci:v00008086d00009D23sv00001028sd000007EAbc0Csc05i00\n")
	if dev == nil {
		t.Fatalf("Failed to parse valid modalias")
	}
	if dev.Class.Name == util.UNKNOWN {
		t.Fatalf("Failed to lookup class name")
	}
}
