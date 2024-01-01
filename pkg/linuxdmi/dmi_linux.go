// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package linuxdmi

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	ghwcontext "github.com/go-hardware/ghw/pkg/context"
	ghwpath "github.com/go-hardware/ghw/pkg/path"
	"github.com/go-hardware/ghw/pkg/util"
)

func Item(ctx context.Context, value string) string {
	paths := ghwpath.New(ctx)
	path := filepath.Join(paths.SysClassDMI, "id", value)

	b, err := os.ReadFile(path)
	if err != nil {
		ghwcontext.Warn(ctx, "Unable to read %s: %s\n", value, err)
		return util.UNKNOWN
	}

	return strings.TrimSpace(string(b))
}
