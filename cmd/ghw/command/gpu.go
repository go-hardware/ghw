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

// gpuCmd represents the gpu command
var gpuCmd = &cobra.Command{
	Use:   "gpu",
	Short: "Show graphics/GPU information for the host system",
	RunE:  showGPU,
}

// showGPU show graphics/GPU information for the host system.
func showGPU(cmd *cobra.Command, args []string) error {
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
	gpu, err := ghw.GPU(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting GPU info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", gpu)

		for _, card := range gpu.GraphicsCards {
			fmt.Printf(" %v\n", card)
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", gpu.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", gpu.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(gpuCmd)
}
