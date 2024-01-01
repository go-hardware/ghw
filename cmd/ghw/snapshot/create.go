//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// staticGlobs is a slice of glob patterns which represent the pseudofiles
	// ghw cares about, and which are independent from host specific topology
	// or configuration, thus are safely represented by a static slice - e.g.
	// they don't need to be discovered at runtime.
	staticGlobs = []string{
		"/proc/cpuinfo",
		"/proc/meminfo",
		"/proc/self/mounts",
		"/sys/devices/system/cpu/cpu*/cache/index*/*",
		"/sys/devices/system/cpu/cpu*/topology/*",
		"/sys/devices/system/memory/block_size_bytes",
		"/sys/devices/system/memory/memory*/online",
		"/sys/devices/system/memory/memory*/state",
		"/sys/devices/system/node/has_*",
		"/sys/devices/system/node/online",
		"/sys/devices/system/node/possible",
		"/sys/devices/system/node/node*/cpu*",
		"/sys/devices/system/node/node*/distance",
		"/sys/devices/system/node/node*/meminfo",
		"/sys/devices/system/node/node*/memory*",
		"/sys/devices/system/node/node*/hugepages/hugepages-*/*",
	}
	// createPaths are always created as-is inside the snapshot filesystem
	createPaths = []string{
		"/sys/block",
	}
)

// New creates a tarball snapshot by copying all the pseudofiles from which ghw
// reads system information into `buildPath` and tarring it up into `outPath`.
func New(
	buildPath string,
	outPath string,
) error {
	snap := &snapshotter{
		buildPath: strings.TrimSuffix(buildPath, string(os.PathSeparator)),
		outPath:   outPath,
	}

	for _, p := range createPaths {
		if err := snap.addDir(p); err != nil {
			return err
		}
	}

	if err := snap.createBlockDevices(); err != nil {
		return err
	}
	if err := snap.copyFileGlobs(staticGlobs); err != nil {
		return err
	}
	if err := snap.copyFileGlobs(netGlobs()); err != nil {
		return err
	}
	if err := snap.copyFileGlobs(pciGlobs()); err != nil {
		return err
	}
	if err := snap.copyFileGlobs(gpuGlobs()); err != nil {
		return err
	}
	return snap.createSnapshot()
}

type snapshotter struct {
	// buildPath is the filepath to the root build directory
	buildPath string
	// outPath is the filepath to the tarball archive produced by the
	// snapshotter
	outPath string
}

// addDir adds a new directory to the snapshot's build directory
func (s *snapshotter) addDir(dir string) error {
	fp := filepath.Join(s.buildPath, dir)
	return os.MkdirAll(fp, os.ModePerm)
}

// createArchive creates the tarball archive.
func (s *snapshotter) createSnapshot() error {
	var f *os.File
	var err error

	if _, err = os.Stat(s.outPath); errors.Is(err, os.ErrNotExist) {
		if f, err = os.Create(s.outPath); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		f, err := os.OpenFile(s.outPath, os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		fs, err := f.Stat()
		if err != nil {
			return err
		}
		if fs.Size() > 0 {
			return fmt.Errorf(
				"file %s already exists and is of size >0",
				s.outPath,
			)
		}
	}
	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()
	return filepath.Walk(s.buildPath, func(path string, fi os.FileInfo, _ error) error {
		if path == s.buildPath {
			return nil
		}
		var link string
		var err error

		if fi.Mode()&os.ModeSymlink != 0 {
			link, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(fi, link)
		if err != nil {
			return err
		}
		hdr.Name = strings.TrimPrefix(
			strings.TrimPrefix(path, s.buildPath),
			string(os.PathSeparator),
		)

		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeReg:
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err = io.Copy(tw, f); err != nil {
				return err
			}
			f.Close()
		}
		return nil
	})
}

// copyFileGlobs copies all the given glob files specs into the tarball's
// build directory, preserving the directory structure. This means you can
// provide a deeply nested filespec like:
//
// - /some/deeply/nested/file*
// and you DO NOT need to build the tree incrementally like
// - /some/
// - /some/deeply/
// ...
// all glob patterns supported in `filepath.Glob` are supported.
func (s *snapshotter) copyFileGlobs(globs []string) error {
	for _, glob := range globs {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return err
		}
		for _, path := range matches {
			if err := s.copyFile(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// Attempting to tar up pseudofiles like /proc/cpuinfo is an exercise in
// futility. Notably, the pseudofiles, when read by syscalls, do not return the
// number of bytes read. This causes the tar writer to write zero-length files.
//
// Instead, it is necessary to build a directory structure in a tmpdir and
// create actual files with copies of the pseudofile contents
func (s *snapshotter) copyFile(path string) error {
	baseDir := filepath.Dir(path)
	if err := s.addDir(baseDir); err != nil {
		return err
	}

	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}
	destPath := filepath.Join(s.buildPath, path)
	if fi.IsDir() {
		// directories must be listed explicitly and created separately.
		if strings.Contains(destPath, "drivers") {
			// When creating snapshots, empty directories are most often
			// useless (but also harmless). Because of this, directories
			// are only created as side effect of copying the files which
			// are inside, and thus directories are never empty. The only
			// notable exception are device driver on linux: in this case,
			// for a number of technical/historical reasons, we care about
			// the directory name, but not about the files which are
			// inside.  Hence, this is the only case on which ghw clones
			// empty directories.
			if err := s.addDir(destPath); err != nil {
				return err
			}
		}
	} else if fi.Mode()&os.ModeSymlink != 0 {
		if err := copyLink(path, destPath); err != nil {
			return err
		}
	} else {
		if err := copyPseudoFile(path, destPath); err != nil && !errors.Is(err, os.ErrPermission) {
			return err
		}
	}
	return nil
}

func copyLink(path, targetPath string) error {
	target, err := os.Readlink(path)
	if err != nil {
		return err
	}
	if err := os.Symlink(target, targetPath); err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return err
	}

	return nil
}

func copyPseudoFile(path, targetPath string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	if _, err = f.Write(buf); err != nil {
		return err
	}
	f.Close()
	return nil
}

type filterFunc func(string) bool

// filterNone allows all content, filtering out none of it
func filterNone(_ string) bool {
	return true
}

// cloneContentByClass copies all the content related to a given device class
// (devClass), possibly filtering out devices whose name does NOT pass a
// filter (filterName). Each entry in `/sys/class/$CLASS` is actually a
// symbolic link. We can filter out entries depending on the link target.
// Each filter is a simple function which takes the entry name or the link
// target and must return true if the entry should be collected, false
// otherwise. Last, explicitly collect a list of attributes for each entry,
// given as list of glob patterns as `subEntries`.
// Return the final list of glob patterns to be collected.
func cloneContentByClass(
	devClass string,
	subEntries []string,
	filterName filterFunc,
	filterLink filterFunc,
) []string {
	var fileSpecs []string

	// warning: don't use the context package here, this means not even the linuxpath package.
	// TODO(fromani) remove the path duplication
	sysClass := filepath.Join("sys", "class", devClass)
	entries, err := os.ReadDir(sysClass)
	if err != nil {
		// we should not import context, hence we can't Warn()
		return fileSpecs
	}
	for _, entry := range entries {
		devName := entry.Name()

		if !filterName(devName) {
			continue
		}

		devPath := filepath.Join(sysClass, devName)
		dest, err := os.Readlink(devPath)
		if err != nil {
			continue
		}

		if !filterLink(dest) {
			continue
		}

		// so, first copy the symlink itself
		fileSpecs = append(fileSpecs, devPath)
		// now we have to clone the content of the actual entry
		// related (and found into a subdir of) the backing hardware
		// device
		devData := filepath.Clean(filepath.Join(sysClass, dest))
		for _, subEntry := range subEntries {
			fileSpecs = append(fileSpecs, filepath.Join(devData, subEntry))
		}
	}

	return fileSpecs
}

// netGlobs returns a slice of strings pertaining to the network interfaces ghw
// cares about. We cannot use a static list because we want to filter away the
// virtual devices, which  ghw doesn't concern itself about. So we need to do
// some runtime discovery.  Additionally, we want to make sure to clone the
// backing device data.
func netGlobs() []string {
	// intentionally avoid to cloning "address" to avoid leaking
	// host-idenfifiable data.
	ifaceEntries := []string{
		"addr_assign_type",
	}

	filterLink := func(linkDest string) bool {
		return !strings.Contains(linkDest, "devices/virtual/net")
	}

	return cloneContentByClass("net", ifaceEntries, filterNone, filterLink)
}

// gpuGlobs returns a slice of strings pertaining to the GPU devices ghw cares
// about. We cannot use a static list because we want to grab only the first
// cardX data (see comment in pkg/gpu/gpu_linux.go) Additionally, we want to
// make sure to clone the backing device data.
func gpuGlobs() []string {
	cardEntries := []string{
		"device",
	}

	filterName := func(cardName string) bool {
		if !strings.HasPrefix(cardName, "card") {
			return false
		}
		if strings.ContainsRune(cardName, '-') {
			return false
		}
		return true
	}

	return cloneContentByClass("drm", cardEntries, filterName, filterNone)
}
