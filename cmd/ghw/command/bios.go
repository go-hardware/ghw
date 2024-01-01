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

// biosCmd represents the bios command
var biosCmd = &cobra.Command{
	Use:   "bios",
	Short: "Show BIOS information for the host system",
	RunE:  showBIOS,
}

// showBIOS shows BIOS host system.
func showBIOS(cmd *cobra.Command, args []string) error {
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
	bios, err := ghw.BIOS(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting BIOS info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", bios)
	case outputFormatJSON:
		fmt.Printf("%s\n", bios.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", bios.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(biosCmd)
}
