package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yaircamilo05/email_to_json/models"
)

var (
	baseURL         string
	username        string
	password        string
	apiKeyIngestion string
	contentType     = "application/json"
)

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error cargando el archivo .env:", err)
	}

	baseURL = os.Getenv("BASEURL")
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
	apiKeyIngestion = os.Getenv("API_KEY_INGESTION")

	if baseURL == "" || username == "" || password == "" || apiKeyIngestion == "" {
		log.Fatal("Faltan variables de entorno necesarias")
	}
}

const (
	errCreatingHTTPRequest = "error creando la petición HTTP: %v"
	errSendingHTTPRequest  = "error enviando la petición HTTP: %v"
	errDecodingResponse    = "error decodificando la respuesta JSON: %v"
	errHTTPResponse        = "error en la respuesta HTTP: %s"
	respSuccess            = "Petición HTTP enviada exitosamente: %s"
)

func GetEmails(querySQL models.Query) (models.SearchResponse, error) {
	url := fmt.Sprintf("%s/_search", baseURL)

	body := models.SearchEmailsResponse{
		SQL:        querySQL,
		Searchtype: "ui",
		Timeout:    0,
	}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		return models.SearchResponse{}, fmt.Errorf(errCreatingHTTPRequest, err)
	}

	fmt.Println(string(jsonBody))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return models.SearchResponse{}, fmt.Errorf(errCreatingHTTPRequest, err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.SearchResponse{}, fmt.Errorf(errSendingHTTPRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.SearchResponse{}, fmt.Errorf(errHTTPResponse, resp.Status)
	}

	var result models.SearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.SearchResponse{}, fmt.Errorf(errDecodingResponse, err)
	}

	fmt.Printf(respSuccess, resp.Status)
	return result, nil
}

func IndexationEmails(data []byte, streamName string) error {
	url := fmt.Sprintf("%s/%s/_json", baseURL, streamName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf(errCreatingHTTPRequest, err)
	}
	req.SetBasicAuth(username, apiKeyIngestion)
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

func GetSchemaByName(SchemaName string) (models.Schema, error) {
	ListSchema, err := GetSchemas()
	if err != nil {
		return models.Schema{}, err
	}
	for _, schema := range ListSchema.List {
		if schema.Name == SchemaName {
			return schema, nil
		}
	}
	return models.Schema{}, fmt.Errorf("schema %s not found", SchemaName)
}
