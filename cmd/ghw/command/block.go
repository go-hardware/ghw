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

// blockCmd represents the block command
var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Show block storage information for the host system",
	RunE:  showBlock,
}

// showBlock show block storage information for the host system.
func showBlock(cmd *cobra.Command, args []string) error {
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
	block, err := ghw.Block(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting block device info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", block)

		for _, disk := range block.Disks {
			fmt.Printf(" %v\n", disk)
			for _, part := range disk.Partitions {
				fmt.Printf("  %v\n", part)
			}
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", block.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", block.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(blockCmd)
}
