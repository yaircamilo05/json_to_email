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

const (
	errCreatingHTTPRequest = "error creando la petición HTTP: %v"
	errSendingHTTPRequest  = "error enviando la petición HTTP: %v"
	errDecodingResponse    = "error decodificando la respuesta JSON: %v"
	errHTTPResponse        = "error en la respuesta HTTP: %s"
	respSuccess            = "Petición HTTP enviada exitosamente: %s"
)

func GetEmails() ([]models.Email, error) {
	url := fmt.Sprintf("%s/_search", baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(errCreatingHTTPRequest, err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errSendingHTTPRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errHTTPResponse, resp.Status)
	}

	var emails []models.Email
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, fmt.Errorf(errDecodingResponse, err)
	}

	fmt.Printf(respSuccess, resp.Status)
	return emails, nil
}

func IndexationEmails(data []byte, streamName string) error {
	url := fmt.Sprintf("%s/%s/_json", baseURL, streamName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf(errCreatingHTTPRequest, err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(errSendingHTTPRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(errHTTPResponse, resp.Status)
	}

	fmt.Printf(respSuccess, resp.Status)
	return nil
}

func GetSchemas() (models.GetSchemasResponse, error) {
	url := fmt.Sprintf("%s/streams?fetchSchema=false", baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.GetSchemasResponse{}, fmt.Errorf(errCreatingHTTPRequest, err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.GetSchemasResponse{}, fmt.Errorf(errSendingHTTPRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.GetSchemasResponse{}, fmt.Errorf(errHTTPResponse, resp.Status)
	}

	var result models.GetSchemasResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.GetSchemasResponse{}, fmt.Errorf(errDecodingResponse, err)
	}

	fmt.Printf(respSuccess, resp.Status)
	return result, nil
}
