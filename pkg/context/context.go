//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package context

import (
	"context"
	"fmt"
	"os"
)

type contextKey string

var (
	optsKey = contextKey("ghw.opts")
)

// ContextModifier sets some value on the context
type ContextModifier func(context.Context) context.Context

// New returns a new context.Context set up with any supplied options/modifiers
func New(mods ...ContextModifier) context.Context {
	ctx := context.TODO()
	for _, mod := range mods {
		ctx = mod(ctx)
	}
	return ctx
}

// OptionsFromContext returns a context's Options struct. If no Options struct
// is contained in the context, the default Options struct is returned.
func OptionsFromContext(ctx context.Context) *Options {
	if ctx == nil {
		return defaultOpts()
	}
	if v := ctx.Value(optsKey); v != nil {
		return v.(*Options)
	}
	return defaultOpts()
}

// WithOptions sets a context's Options. If the supplied Options' fields are
// nil, the default Options derived from environs variables are used to
// populated the Options.
func WithOptions(opts *Options) ContextModifier {
	return func(ctx context.Context) context.Context {
		// merge the supplied options with the default options derived from
		// environs variables
		defOpts := defaultOpts()
		if opts.RootMountpoint == nil {
			opts.RootMountpoint = defOpts.RootMountpoint
		}
		if opts.DisableWarnings == nil {
			opts.DisableWarnings = defOpts.DisableWarnings
		}
		if opts.DisableExternalTools == nil {
			opts.DisableExternalTools = defOpts.DisableExternalTools
		}
		return context.WithValue(ctx, optsKey, opts)
	}
}

// WithDisableWarnings disables warning messages (mostly about permissions
// issues ghw encounters trying to discovery hardware information).
func WithDisableWarnings() ContextModifier {
	return func(ctx context.Context) context.Context {
		opts := OptionsFromContext(ctx)
		_true := true
		opts.DisableWarnings = &_true
		return context.WithValue(ctx, optsKey, opts)
	}
}

// WithDisableExternalTools disables warning messages (mostly about permissions
// issues ghw encounters trying to discovery hardware information).
func WithDisableExternalTools() ContextModifier {
	return func(ctx context.Context) context.Context {
		opts := OptionsFromContext(ctx)
		_true := true
		opts.DisableExternalTools = &_true
		return context.WithValue(ctx, optsKey, opts)
	}
}

// WithRootMountpoint sets the root mountpoint ghw uses when querying system
// information.
func WithRootMountpoint(path string) ContextModifier {
	return func(ctx context.Context) context.Context {
		opts := OptionsFromContext(ctx)
		opts.RootMountpoint = &path
		return context.WithValue(ctx, optsKey, opts)
	}
}

// WithPathOverrides supplies path-specific overrides for the context
func WithPathOverrides(overrides PathOverrides) ContextModifier {
	return func(ctx context.Context) context.Context {
		opts := OptionsFromContext(ctx)
		opts.PathOverrides = &overrides
		return context.WithValue(ctx, optsKey, opts)
	}
}

// Warn prints a warning message. If the DisableWarnings Option has been set
// (or the corresponding GHW_DISABLE_WARNINGS environs variable is true,
// nothing is printed.
func Warn(ctx context.Context, msg string, args ...interface{}) {
	opts := OptionsFromContext(ctx)
	if opts.DisableWarnings == nil || !*opts.DisableWarnings {
		fmt.Fprintf(os.Stderr, "WARNING: "+msg, args...)
	}
}
