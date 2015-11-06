package utility

import (
	"io"
	"os"
)

func Cp(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func Mv(src, dst string) error {
	_, err := os.Open(src)
	if err != nil {
		return err
	}

	err = Cp(src, dst)
	if err != nil {
		return err
	}

	err = os.Remove(src)

	return err
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// Checks whether a file exists at the given filepath
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)

	return err == nil
}
