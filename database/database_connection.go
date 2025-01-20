package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yaircamilo05/email_to_json/models"
)

const (
	baseURL     = "http://localhost:5080/api/prueba"
	username    = "root@example.com"
	password    = "T9uplVBu16xjKUrd"
	contentType = "application/json"
)

func GetEmails() ([]models.Email, error) {
	url := fmt.Sprintf("%s/_search", baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando la petición HTTP: %v", err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando la petición HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error en la respuesta HTTP: %s", resp.Status)
	}

	var emails []models.Email
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, fmt.Errorf("error decodificando la respuesta JSON: %v", err)
	}

	fmt.Printf("Petición HTTP enviada exitosamente: %s\n", resp.Status)
	return emails, nil
}

func IndexationEmails(data []byte, streamName string) error {
	url := fmt.Sprintf("%s/%s/_json", baseURL, streamName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creando la petición HTTP: %v", err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando la petición HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta HTTP: %s", resp.Status)
	}

	fmt.Printf("Petición HTTP enviada exitosamente: %s\n", resp.Status)
	return nil
}
