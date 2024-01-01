//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/go-hardware/ghw"
	"github.com/go-hardware/ghw/cmd/ghw/snapshot"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// cpuCmd represents the cpu command
var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Show CPU information for the host system",
	RunE:  showCPU,
}

// showCPU show CPU information for the host system.
func showCPU(cmd *cobra.Command, args []string) error {
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
	cpu, err := ghw.CPU(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting CPU info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", cpu)

		for _, proc := range cpu.Processors {
			fmt.Printf(" %v\n", proc)
			for _, core := range proc.Cores {
				fmt.Printf("  %v\n", core)
			}
			if len(proc.Capabilities) > 0 {
				// pretty-print the (large) block of capability strings into rows
				// of 6 capability strings
				rows := int(math.Ceil(float64(len(proc.Capabilities)) / float64(6)))
				for row := 1; row < rows; row = row + 1 {
					rowStart := (row * 6) - 1
					rowEnd := int(math.Min(float64(rowStart+6), float64(len(proc.Capabilities))))
					rowElems := proc.Capabilities[rowStart:rowEnd]
					capStr := strings.Join(rowElems, " ")
					if row == 1 {
						fmt.Printf("  capabilities: [%s\n", capStr)
					} else if rowEnd < len(proc.Capabilities) {
						fmt.Printf("                 %s\n", capStr)
					} else {
						fmt.Printf("                 %s]\n", capStr)
					}
				}
			}
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", cpu.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", cpu.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(cpuCmd)
}
