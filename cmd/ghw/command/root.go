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

const (
	appName   = "ghw"
	shortDesc = "ghw - discover hardware information and query system resource usage"
	longDesc  = `
          __
 .-----. |  |--. .--.--.--.
 |  _  | |     | |  |  |  |
 |___  | |__|__| |________|
 |_____|

Discover hardware information and query system resource usage.

https://github.com/go-hardware/ghw
`
	outputFormatHuman = "human"
	outputFormatJSON  = "json"
	outputFormatYAML  = "yaml"
	usageOutputFormat = `Output format.
Choices are 'json','yaml', and 'human'.`
	usageSnapshotPath = `Snapshot path.
If you want ghw to examine a snapshot file, pass the path to the snapshot.`
)

var (
	version       string
	buildHash     string
	buildDate     string
	debug         bool
	outputFormat  string
	outputFormats = []string{
		outputFormatHuman,
		outputFormatJSON,
		outputFormatYAML,
	}
	pretty             bool
	snapshotPath       string
	snapshotExpandPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: shortDesc,
	Args:  validateRootCommand,
	Long:  longDesc,
	RunE:  showAll,
}

func showAll(cmd *cobra.Command, args []string) error {
	var err error
	if snapshotPath != "" {
		snapshotExpandPath, err = os.MkdirTemp("", "ghw-snapshot")
		if err != nil {
			return err
		}
		defer os.RemoveAll(snapshotExpandPath)
		if err = snapshot.Expand(snapshotPath, snapshotExpandPath); err != nil {
			return err
		}
	}

	switch outputFormat {
	case outputFormatHuman:
		if err := showBlock(cmd, args); err != nil {
			return err
		}
		if err := showCPU(cmd, args); err != nil {
			return err
		}
		if err := showGPU(cmd, args); err != nil {
			return err
		}
		if err := showMemory(cmd, args); err != nil {
			return err
		}
		if err := showNetwork(cmd, args); err != nil {
			return err
		}
		if err := showTopology(cmd, args); err != nil {
			return err
		}
		if err := showChassis(cmd, args); err != nil {
			return err
		}
		if err := showBIOS(cmd, args); err != nil {
			return err
		}
		if err := showBaseboard(cmd, args); err != nil {
			return err
		}
		if err := showProduct(cmd, args); err != nil {
			return err
		}
	case outputFormatJSON:
		system, err := ghw.System(context.TODO())
		if err != nil {
			return errors.Wrap(err, "error getting system info")
		}
		fmt.Printf("%s\n", system.JSONString(pretty))
	case outputFormatYAML:
		system, err := ghw.System(context.TODO())
		if err != nil {
			return errors.Wrap(err, "error getting system info")
		}
		fmt.Printf("%s", system.YAMLString())
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute(v string, bh string, bd string) {
	version = v
	buildHash = bh
	buildDate = bd

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func haveValidOutputFormat() bool {
	for _, choice := range outputFormats {
		if choice == outputFormat {
			return true
		}
	}
	return false
}

// validateRootCommand ensures any CLI options or arguments are valid,
// returning an error if not
func validateRootCommand(rootCmd *cobra.Command, args []string) error {
	if !haveValidOutputFormat() {
		return fmt.Errorf("invalid output format %q", outputFormat)
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(
		&debug, "debug", false, "Enable or disable debug mode",
	)
	rootCmd.PersistentFlags().StringVarP(
		&outputFormat,
		"format", "f",
		outputFormatHuman,
		usageOutputFormat,
	)
	rootCmd.PersistentFlags().StringVarP(
		&snapshotPath,
		"snapshot", "s",
		"",
		usageSnapshotPath,
	)
	rootCmd.PersistentFlags().BoolVar(
		&pretty, "pretty", false, "When outputting JSON, use indentation",
	)
}
