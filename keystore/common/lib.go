package common

import (
	"fmt"
	"os"
)

func WriteFile(filepath string, b []byte, mode uint32) error {

	tmp := fmt.Sprintf("%s.tmp", filepath)
	unlink(tmp)

	fd, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("can't create key file %s: %s", tmp, err)
	}

	_, err = fd.Write(b)
	if err != nil {
		fd.Close()
		return fmt.Errorf("can't write %v bytes to %s: %s", len(b), tmp, err)
	}

	fd.Close() // we ignore close(2) errors; unrecoverable anyway.

	os.Rename(tmp, filepath)
	return nil
}

func unlink(f string) error {
	st, err := os.Stat(f)
	if err == nil {
		if !st.Mode().IsRegular() {
			return fmt.Errorf("%s can't be unlinked. Not a regular file?", f)
		}

		os.Remove(f)
		return nil
	}

	return err
}
