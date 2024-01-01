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

// chassisCmd represents the chassis command
var chassisCmd = &cobra.Command{
	Use:   "chassis",
	Short: "Show chassis information for the host system",
	RunE:  showChassis,
}

// showChassis shows chassis information for the host system.
func showChassis(cmd *cobra.Command, args []string) error {
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
	chassis, err := ghw.Chassis(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting chassis info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", chassis)
	case outputFormatJSON:
		fmt.Printf("%s\n", chassis.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", chassis.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(chassisCmd)
}
