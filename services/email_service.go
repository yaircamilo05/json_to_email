package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yaircamilo05/email_to_json/database"
	"github.com/yaircamilo05/email_to_json/models"
	"github.com/yaircamilo05/email_to_json/utils"
)

func ProcessEmails(dirPath, streamName string) error {
	directories := []string{}
	err := utils.FindDirectories(dirPath, &directories)
	if err != nil {
		return err
	}

	directoriesWithFiles := utils.FilterDirectoriesWithFiles(directories)
	var wg sync.WaitGroup

	for _, dir := range directoriesWithFiles {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			dirEmails, err := processDirectory(dir)
			if err != nil {
				fmt.Println("Error al procesar el directorio:", err)
				return
			}

			jsonData, err := json.Marshal(dirEmails)
			if err != nil {
				fmt.Println("Error al convertir a JSON:", err)
				return
			}

			err = database.IndexationEmails(jsonData, streamName)
			if err != nil {
				if strings.Contains(err.Error(), "413 Request Entity Too Large") {
					handleLargeRequest(dirEmails, streamName)
				} else {
					fmt.Println("Error al enviar la petición HTTP:", err)
				}
			}
		}(dir)
	}

	wg.Wait()
	return nil
}

func processDirectory(path string) ([]models.Email, error) {
	emails := []models.Email{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error al leer el directorio: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(path, entry.Name())
			email, err := models.ReadEmail(filePath)
			if err != nil {
				fmt.Println("Error al leer el archivo:", err)
				continue
			}
			emails = append(emails, email)
		}
	}

	return emails, nil
}

func handleLargeRequest(emails []models.Email, streamName string) {
	// Dividir el arreglo de emails en dos partes y reenviar
	mid := len(emails) / 2
	batch1 := emails[:mid]
	batch2 := emails[mid:]

	jsonData1, err := json.Marshal(batch1)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	jsonData2, err := json.Marshal(batch2)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	err = database.IndexationEmails(jsonData1, streamName)
	if err != nil {
		fmt.Println("Error al enviar la petición HTTP:", err)
	}
	fmt.Println("Petición HTTP enviada exitosamente en el segundo intento")

	err = database.IndexationEmails(jsonData2, streamName)
	if err != nil {
		fmt.Println("Error al enviar la petición HTTP:", err)
	}
	fmt.Println("Petición HTTP enviada exitosamente en el segundo intento")
}

func GetEmails(req models.GetAllEmailsRequest) (models.SearchResponse, error) {
	schemaSelected, err := database.GetSchemaByName(req.Schema)
	if err != nil {
		return models.SearchResponse{}, fmt.Errorf("error al obtener el schema: %v", err)
	}

	querySQL := models.Query{
		SQL:       req.SQL,
		From:      req.From,
		Size:      req.Size,
		StartTime: schemaSelected.Stats.DocTimeMin,
		EndTime:   schemaSelected.Stats.DocTimeMax,
	}

	listEmails, err := database.GetEmails(querySQL)
	if err != nil {
		return models.SearchResponse{}, fmt.Errorf("error al obtener los emails: %v", err)
	}

	return listEmails, nil
}
