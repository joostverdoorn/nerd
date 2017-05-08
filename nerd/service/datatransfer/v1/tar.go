package v1datatransfer

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dchest/safefile"
	"github.com/pkg/errors"
)

//tardir archives the given directory and writes bytes to w.
func tardir(dir string, w io.Writer) (err error) {
	tw := tar.NewWriter(w)
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return errors.Wrapf(err, "failed to determine path '%s' relative to '%s'", path, dir)
		}

		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return errors.Wrapf(err, "failed to open file '%s'", rel)
		}

		err = tw.WriteHeader(&tar.Header{
			Name:    rel,
			Mode:    int64(fi.Mode()),
			ModTime: fi.ModTime(),
			Size:    fi.Size(),
		})
		if err != nil {
			if err == io.ErrClosedPipe {
				return err
			}
			return errors.Wrapf(err, "failed to write tar header for '%s'", rel)
		}

		n, err := io.Copy(tw, f)
		// fmt.Printf("%v %v\n", path, n)
		if err != nil {
			if err == io.ErrClosedPipe {
				return err
			}
			return errors.Wrapf(err, "failed to write tar file for '%s'", rel)
		}

		if n != fi.Size() {
			return errors.Errorf("unexpected nr of bytes written to tar, saw '%d' on-disk but only wrote '%d', is directory '%s' in use?", fi.Size(), n, dir)
		}

		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "failed to walk dir '%s'", dir)
	}

	if err = tw.Close(); err != nil {
		if err == io.ErrClosedPipe {
			return err
		}
		return errors.Wrap(err, "failed to write remaining data")
	}

	return nil
}

//untardir untars an archive from the reader to a directory on disk.
func untardir(dir string, r io.Reader) (err error) {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return errors.Wrap(err, "failed to read next tar header")
		}

		path := filepath.Join(dir, hdr.Name)
		err = os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return errors.Wrap(err, "failed to create dirs")
		}

		f, err := safefile.Create(path, os.FileMode(hdr.Mode))
		if err != nil {
			return errors.Wrap(err, "failed to create tmp safe file")
		}

		defer f.Close()
		n, err := io.Copy(f, tr)
		if err != nil {
			return errors.Wrap(err, "failed to write file content to tmp file")
		}

		if n != hdr.Size {
			return errors.Errorf("unexpected nr of bytes written, wrote '%d' saw '%d' in tar hdr", n, hdr.Size)
		}

		err = f.Commit()
		if err != nil {
			return errors.Wrap(err, "failed to swap old file for tmp file")
		}

		err = os.Chtimes(path, time.Now(), hdr.ModTime)
		if err != nil {
			return errors.Wrap(err, "failed to change times of tmp file")
		}
	}

	return nil
}
