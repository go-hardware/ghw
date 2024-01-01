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

// netCmd represents the net command
var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Show network information for the host system",
	RunE:  showNetwork,
}

// showNetwork show network information for the host system.
func showNetwork(cmd *cobra.Command, args []string) error {
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
	net, err := ghw.Network(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting network info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", net)

		for _, nic := range net.NICs {
			fmt.Printf(" %v\n", nic)

			enabledCaps := make([]int, 0)
			for x, cap := range nic.Capabilities {
				if cap.IsEnabled {
					enabledCaps = append(enabledCaps, x)
				}
			}
			if len(enabledCaps) > 0 {
				fmt.Printf("  enabled capabilities:\n")
				for _, x := range enabledCaps {
					fmt.Printf("   - %s\n", nic.Capabilities[x].Name)
				}
			}
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", net.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", net.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(netCmd)
}
