package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Define el directorio que quieres explorar
	dir := "C:/Users/yairc/OneDrive/Documentos/universidad/semestre 7/truora 2/json_separados/"
	// Define la función que ejecutará el comando para un archivo
	executeCommand := func(file string) {
		cmd := exec.Command("curl", "-u", "yaircamilo05@gmail.com:6B4OG0158f27SZ3Av9Cr", "-k", "https://api.openobserve.ai/api/yair_camilo_organization_15016_FGL9xXNbnwFSWC2/emailsA/_json", "-d", "@"+file)
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error ejecutando el comando para el archivo %s: %s\n", file, err)
			return
		}
		fmt.Printf("Comando ejecutado exitosamente para el archivo %s\n", file)
		fmt.Printf("Salida del comando:\n%s\n", output)
	}

	// Recorre el directorio y sus subdirectorios
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error al acceder a la ruta:", path, err)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Si el archivo es un archivo JSON, ejecuta el comando
			executeCommand(path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error al recorrer el directorio:", err)
	}
}