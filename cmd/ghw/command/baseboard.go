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

// baseboardCmd represents the baseboard command
var baseboardCmd = &cobra.Command{
	Use:   "baseboard",
	Short: "Show baseboard information for the host system",
	RunE:  showBaseboard,
}

// showBaseboard shows baseboard information for the host system.
func showBaseboard(cmd *cobra.Command, args []string) error {
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
	baseboard, err := ghw.Baseboard(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting baseboard info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", baseboard)
	case outputFormatJSON:
		fmt.Printf("%s\n", baseboard.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", baseboard.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(baseboardCmd)
}
