package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindDirectories(path string, directories *[]string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		subPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			*directories = append(*directories, subPath)
			if err := FindDirectories(subPath, directories); err != nil {
				return err
			}
		}
	}
	return nil
}

func FilterDirectoriesWithFiles(directories []string) []string {
	var result []string
	for _, dir := range directories {
		if hasFiles(dir) {
			result = append(result, dir)
		}
	}
	return result
}

func hasFiles(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error al leer el directorio:", err)
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			return true
		}
	}
	return false
}
