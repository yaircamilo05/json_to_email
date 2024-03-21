package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Estructura de ejemplo para los datos JSON
type Email struct {
	MessageID               string `json:"Message-ID"`
	Date                    string `json:"Date"`
	From                    string `json:"From"`
	To                      string `json:"To"`
	Subject                 string `json:"Subject"`
	MimeVersion             string `json:"Mime-Version"`
	ContentType             string `json:"Content-Type"`
	ContentTransferEncoding string `json:"Content-Transfer-Encoding"`
	XFrom                   string `json:"X-From"`
	XTo                     string `json:"X-To"`
	XCC                     string `json:"X-cc"`
	XBCC                    string `json:"X-bcc"`
	XFolder                 string `json:"X-Folder"`
	XOrigin                 string `json:"X-Origin"`
	XFileName               string `json:"X-FileName"`
	Body                    string `json:"Body"`
}

func main() {
	// Obtener la carpeta de entrada del argumento de línea de comandos

	inputFolder := "C:/Users/yairc/OneDrive/Documentos/universidad/semestre 7/truora 2/json_separados/separados/todavia muy grandes"

	// Lee todos los archivos JSON en la carpeta de entrada
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		fmt.Println("Error al leer la carpeta:", err)
		return
	}

	// Define el tamaño máximo de cada archivo
	maxSize := 5 * 1024 * 1024 // Tamaño máximo en bytes

	// Carpeta de salida para los archivos partidos
	outputFolder := "C:/Users/yairc/OneDrive/Documentos/universidad/semestre 7/truora 2/json_separados/separados/todavia muy grandes/slice5MB"

	// Recorre cada archivo JSON en la carpeta
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			// Lee el archivo JSON
			jsonFile, err := os.Open(filepath.Join(inputFolder, file.Name()))
			if err != nil {
				fmt.Println("Error al abrir el archivo JSON:", err)
				continue
			}
			defer jsonFile.Close()

			// Lee el contenido del archivo
			byteValue, _ := ioutil.ReadAll(jsonFile)

			// Decodifica el JSON
			var emails []Email
			err = json.Unmarshal(byteValue, &emails)
			if err != nil {
				fmt.Println("Error al decodificar el JSON:", err)
				continue
			}

			// Variables para el archivo de salida
			outputPrefix := filepath.Join(outputFolder, file.Name()[:len(file.Name())-5]) // Elimina la extensión .json del nombre del archivo
			currentChunk := 1
			currentSize := 0
			currentData := []Email{}

			// Recorre los datos y escribe en múltiples archivos
			for _, email := range emails {
				// Codifica el objeto actual
				jsonData, err := json.Marshal(email)
				if err != nil {
					fmt.Println("Error al codificar el objeto JSON:", err)
					continue
				}

				// Calcula el tamaño del objeto
				objSize := len(jsonData)

				// Si agregar este objeto supera el tamaño máximo, guarda el archivo actual y comienza uno nuevo
				if currentSize+objSize > maxSize {
					outputFileName := fmt.Sprintf("%s_part%d.json", outputPrefix, currentChunk)
					err := writeToFile(outputFileName, currentData)
					if err != nil {
						fmt.Println("Error al escribir en el archivo:", err)
						return
					}

					// Incrementa el contador de archivos y reinicia las variables
					currentChunk++
					currentSize = 0
					currentData = []Email{}
				}

				// Agrega el objeto actual al conjunto de datos actual
				currentData = append(currentData, email)
				currentSize += objSize
			}

			// Escribe el último archivo si quedan datos
			if len(currentData) > 0 {
				outputFileName := fmt.Sprintf("%s_part%d.json", outputPrefix, currentChunk)
				err := writeToFile(outputFileName, currentData)
				if err != nil {
					fmt.Println("Error al escribir en el archivo:", err)
					return
				}
			}
		}
	}
}

// Función para escribir datos en un archivo
func writeToFile(filename string, data []Email) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Codifica los datos y escribe en el archivo
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	fmt.Println("Archivo creado:", filename)
	return nil
}
