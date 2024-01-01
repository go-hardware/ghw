//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/go-hardware/ghw/cmd/ghw/snapshot"
	"github.com/spf13/cobra"
)

const (
	snapshotLongDesc = `
Create a tarball snapshot containing system information (Linux only)

This snapshot can then be read by ghw. Mostly useful for testing
and debugging ghw.
`
)

var (
	// output filepath to save snapshot to
	outPath string
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Create a snapshot containing system information (Linux only)",
	Long:  snapshotLongDesc,
	RunE:  doSnapshot,
}

func doSnapshot(cmd *cobra.Command, args []string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("ghw snapshot is currently only supported on Linux")
	}
	buildDir, err := os.MkdirTemp("", "ghw-snapshot")
	if err != nil {
		return err
	}
	defer os.RemoveAll(buildDir)

	if err = snapshot.New(buildDir, outPath); err != nil {
		return err
	}
	fmt.Println("successfully wrote snapshot to", outPath)
	return nil
}

func systemFingerprint() (string, error) {
	hn, err := os.Hostname()
	if err != nil {
		return "unknown", err
	}
	m := md5.New()
	_, err = io.WriteString(m, hn)
	if err != nil {
		return "unknown", err
	}
	return fmt.Sprintf("%x", m.Sum(nil)), nil
}

func defaultOutPath() string {
	fp, err := systemFingerprint()
	if err != nil {
		return "unknown"
	}
	return fmt.Sprintf("%s-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH, fp)
}

func init() {
	snapshotCmd.PersistentFlags().StringVarP(
		&outPath,
		"out", "o",
		defaultOutPath(),
		"Output file path. Defaults to file in current directory with name $OS-$ARCH-$HASHSYSTEMNAME.tar.gz",
	)
	rootCmd.AddCommand(snapshotCmd)
}
