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
	emails := []models.Email{}
	var mu sync.Mutex
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
			mu.Lock()
			emails = append(emails, dirEmails...)
			mu.Unlock()
		}(dir)
	}

	wg.Wait()

	jsonData, err := json.Marshal(emails)
	if err != nil {
		return fmt.Errorf("error al convertir a JSON: %v", err)
	}

	err = database.IndexationEmails(jsonData, streamName)
	if err != nil {
		if strings.Contains(err.Error(), "413 Request Entity Too Large") {
			handleLargeRequest(emails, streamName)
		} else {
			return err
		}
	}

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

func GetEmails() ([]models.Email, error) {
	return database.GetEmails()
}
