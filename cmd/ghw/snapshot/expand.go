//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Expand expands the given snapshot into a target directory. If the target
// directory does not exist, this function creates it.
func Expand(
	fromPath string, // the path to the snapshot archive to be expanded
	toPath string, // the directory to expand into
) error {
	if !exists(toPath) {
		if err := os.MkdirAll(toPath, os.ModePerm); err != nil {
			return err
		}
	} else if !isEmptyDir(toPath) {
		return fmt.Errorf("target directory %s is not empty", toPath)
	}
	snap, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer snap.Close()
	return untar(snap, toPath)
}

// Untar extracts data from the given reader (providing data in tar.gz format)
// and unpacks it in the given directory.
func untar(r io.Reader, toPath string) error {
	var err error
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			// we are done
			return nil
		}

		if err != nil {
			// bail out
			return err
		}

		if header == nil {
			// TODO: how come?
			continue
		}

		target := filepath.Join(toPath, header.Name)
		mode := os.FileMode(header.Mode)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(target, mode)
			if err != nil {
				return err
			}

		case tar.TypeReg:
			dst, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, mode)
			if err != nil {
				return err
			}

			_, err = io.Copy(dst, tr)
			if err != nil {
				return err
			}

			dst.Close()

		case tar.TypeSymlink:
			err = os.Symlink(header.Linkname, target)
			if err != nil {
				return err
			}
		}
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func isEmptyDir(name string) bool {
	entries, err := os.ReadDir(name)
	if err != nil {
		return false
	}
	return len(entries) == 0
}
