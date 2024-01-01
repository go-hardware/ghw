//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"context"
	"fmt"
	"os"

	"github.com/go-hardware/ghw"
	"github.com/go-hardware/ghw/cmd/ghw/snapshot"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// pciCmd represents the pci command
var pciCmd = &cobra.Command{
	Use:   "pci",
	Short: "Show information about PCI devices on the host system",
	RunE:  showPCI,
}

// showPCI shows information for PCI devices on the host system.
func showPCI(cmd *cobra.Command, args []string) error {
	var err error
	if snapshotPath != "" && snapshotExpandPath == "" {
		snapshotExpandPath, err = os.MkdirTemp("", "ghw-snapshot")
		if err != nil {
			return err
		}
		defer os.RemoveAll(snapshotExpandPath)
		if err = snapshot.Expand(snapshotPath, snapshotExpandPath); err != nil {
			return err
		}
	}
	ctx := context.TODO()
	if snapshotExpandPath != "" {
		ctx = ghw.NewContext(ghw.WithRootMountpoint(snapshotExpandPath))
	}
	pci, err := ghw.PCI(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting PCI info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", pci)
	case outputFormatJSON:
		fmt.Printf("%s\n", pci.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", pci.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(pciCmd)
}
