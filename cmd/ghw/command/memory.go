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

// memoryCmd represents the memory command
var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Show memory information for the host system",
	RunE:  showMemory,
}

// showMemory show memory information for the host system.
func showMemory(cmd *cobra.Command, args []string) error {
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
	mem, err := ghw.Memory(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting memory info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", mem)
	case outputFormatJSON:
		fmt.Printf("%s\n", mem.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", mem.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(memoryCmd)
}
