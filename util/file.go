package util

import "os"

// EnsureDir verifica se o diretório existe; se não existir, tenta criá-lo.
func EnsureDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return os.MkdirAll(dirName, 0755)
	}
	return nil
}

// WriteTargetsToFile writes the slice of targets (IPs) to the specified file,
// one target per line.
func WriteTargetsToFile(filePath string, targets []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, target := range targets {
		if _, err := file.WriteString(target + "\n"); err != nil {
			return err
		}
	}
	return nil
}
