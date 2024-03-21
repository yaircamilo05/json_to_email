package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define el directorio que contiene los archivos JSON
	dir := "C:/Users/yairc/OneDrive/Documentos/universidad/semestre 7/truora 2/json_separados/"

	// Lee el directorio
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error al leer el directorio:", err)
		return
	}

	// Procesa cada archivo en el directorio
	for _, entry := range entries {
		// Ignora los directorios
		if entry.IsDir() {
			continue
		}

		// Ignora los archivos que no son JSON
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		// Lee el archivo de entrada en memoria
		input, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			fmt.Println("Error al abrir el archivo de entrada:", err)
			continue
		}

		// Convierte el contenido del archivo a una cadena
		content := string(input)

		// Añade '[' al inicio y ']' al final
		content = "[" + content + "]"

		// Reemplaza todas las ocurrencias de '}' con '},'
		content = strings.ReplaceAll(content, "}", "},")

		// Reemplaza la última coma con un espacio vacío
		lastComma := strings.LastIndex(content, ",")
		if lastComma != -1 {
			content = content[:lastComma] + content[lastComma+1:]
		}

		// Escribe el contenido modificado de nuevo al archivo
		err = os.WriteFile(filepath.Join(dir, entry.Name()), []byte(content), 0644)
		if err != nil {
			fmt.Println("Error al escribir en el archivo de entrada:", err)
		}
	}
}