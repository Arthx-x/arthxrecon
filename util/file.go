package util

import "os"

// EnsureDirectory creates the directory if it does not exist.
func EnsureDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}
