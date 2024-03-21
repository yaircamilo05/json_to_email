package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const MaxFileSize = 14000000 // 14MB

func main() {
	// Define el directorio que quieres explorar
	dir := "C:/Users/yairc/OneDrive/Documentos/universidad/semestre 7/truora 2/json_separados/separados/todavia muy grandes/slice5MB/"

	// Recorre el directorio y sus subdirectorios
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error al acceder a la ruta:", path, err)
			return err
		}
		if !info.IsDir() && info.Size() > MaxFileSize {
			// Si el archivo es más grande que el tamaño máximo, elimínalo
			err = os.Remove(path)
			if err != nil {
				fmt.Println("Error al eliminar el archivo:", err)
				return err
			}
			fmt.Printf("Archivo %s eliminado.\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error al recorrer el directorio:", err)
	}
}
