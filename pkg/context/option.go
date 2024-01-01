//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package context

import (
	"github.com/jaypipes/envutil"
)

const (
	defaultRootMountpoint       = "/"
	defaultDisableWarnings      = false
	defaultDisableExternalTools = false
)

const (
	envKeyChroot               = "GHW_CHROOT"
	envKeyRootMountpoint       = "GHW_ROOT_MOUNTPOINT"
	envKeyDisableWarnings      = "GHW_DISABLE_WARNINGS"
	envKeyDisableTools         = "GHW_DISABLE_TOOLS"
	envKeyDisableExternalTools = "GHW_DISABLE_EXTERNAL_TOOLS"
)

// PathOverrides is a map, keyed by the string name of a mount path, of
// override paths
type PathOverrides map[string]string

// Options contains configuration options that control how ghw behaves.
type Options struct {
	// RootMountpoint contains an alternate root mountpoint.
	//
	// To facilitate querying of sysfs filesystems that are bind-mounted to a
	// non-default root mountpoint, we allow users to set the
	// GHW_ROOT_MOUNTPOINT environs variable to an alternate mountpoint. For
	// instance, assume that the user of ghw is a Golang binary being executed
	// from an application container that has certain host filesystems
	// bind-mounted into the container at /host. The user would ensure the
	// GHW_ROOT_MOUNTPOINT environs variable is set to "/host" and ghw will
	// build its paths from that location instead of /.
	//
	// NOTE(jaypipes): If the deprecated GHW_CHROOT environs variable is set,
	// we will use that value. Please use the GHW_ROOT_MOUNTPOINT environs
	// variable instead.
	RootMountpoint *string

	// PathOverrides optionally allows to override the default paths ghw uses
	// internally to learn about the system resources.
	PathOverrides *PathOverrides

	// DisableWarnings tells ghw not to output warnings to stderr.
	//
	// Set the GHW_DISABLE_WARNINGS environs variable to 1 or any truthy value
	// to disable external tools.
	DisableWarnings *bool

	// DisableExternalTools tells ghw to not call any external program to learn
	// about the hardware. The default is to use such tools if available.
	//
	// Set the GHW_DISABLE_EXTERNAL_TOOLS environs variable to 1 or any truthy
	// value to disable external tools.
	//
	// NOTE(jaypipes): If the deprecated GHW_DISABLE_TOOLS environs variable is
	// set, we will use that value. Please use the GHW_DISABLE_EXTERNAL_TOOLS
	// environs variable instead.
	DisableExternalTools *bool
}

// defaultOpts returns the default set of options derived from any environs
// variables.
func defaultOpts() *Options {
	envDefaultRootMountpoint := envutil.WithDefault(
		envKeyRootMountpoint,
		envutil.WithDefault(
			envKeyChroot,
			defaultRootMountpoint,
		),
	)
	envDefaultDisableWarnings := envutil.WithDefaultBool(
		envKeyDisableWarnings,
		defaultDisableWarnings,
	)
	envDefaultDisableExternalTools := envutil.WithDefaultBool(
		envKeyDisableExternalTools,
		envutil.WithDefaultBool(
			envKeyDisableTools,
			defaultDisableExternalTools,
		),
	)
	return &Options{
		RootMountpoint:       &envDefaultRootMountpoint,
		DisableWarnings:      &envDefaultDisableWarnings,
		DisableExternalTools: &envDefaultDisableExternalTools,
	}
}
