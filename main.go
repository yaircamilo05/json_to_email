package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

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
	dirPath := "C:/Users/User/Documents/prueba_truora/enron_mail_20110402/enron_mail_20110402/maildir"
	directories := []string{}
	err := findDirectories(dirPath, &directories)
	if err != nil {
		fmt.Println("Error al buscar carpetas:", err)
		return
	}
	fmt.Println("Directorios encontrados:", directories)
	fmt.Println("Cantidad de directorios encontrados: ", len(directories))

	directoriesWithFiles := filterDirectoriesWithFiles(directories)
	fmt.Println("Directorios con al menos un archivo:", directoriesWithFiles)
	fmt.Println("Cantidad de directorios con archivos:", len(directoriesWithFiles))

	core(directoriesWithFiles)
}

func findDirectories(path string, directories *[]string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		subPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			*directories = append(*directories, subPath)
			if err := findDirectories(subPath, directories); err != nil {
				return err
			}
		}
	}
	return nil
}

func filterDirectoriesWithFiles(directories []string) []string {
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

func core(directories []string) {
	var wg sync.WaitGroup
	for _, dir := range directories {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			processDirectories(dir)
		}(dir)
	}
	wg.Wait()
}

func processDirectories(path string) {
	var wg sync.WaitGroup
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error al leer el directorio:", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			wg.Add(1)
			go func(entry os.DirEntry) {
				defer wg.Done()
				filePath := filepath.Join(path, entry.Name())
				email, err := readEmail(filePath)
				if err != nil {
					fmt.Println("Error al leer el archivo:", err)
					return
				}
				jsonData, err := json.Marshal(email)
				if err != nil {
					fmt.Println("Error al convertir a JSON:", err)
					return
				}
				err = sendHTTPRequest(jsonData)
				if err != nil {
					fmt.Println("Error al enviar la petición HTTP:", err)
				}
			}(entry)
		}
	}
	wg.Wait()
}

func readEmail(filePath string) (Email, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Email{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	email := Email{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") {
			split := strings.SplitN(line, ":", 2)
			field := strings.TrimSpace(split[0])
			value := strings.TrimSpace(split[1])
			switch field {
			case "Message-ID":
				email.MessageID = value
			case "Date":
				email.Date = value
			case "From":
				email.From = value
			case "To":
				email.To = value
			case "Subject":
				email.Subject = value
			case "Mime-Version":
				email.MimeVersion = value
			case "Content-Type":
				email.ContentType = value
			case "Content-Transfer-Encoding":
				email.ContentTransferEncoding = value
			case "X-From":
				email.XFrom = value
			case "X-To":
				email.XTo = value
			case "X-cc":
				email.XCC = value
			case "X-bcc":
				email.XBCC = value
			case "X-Folder":
				email.XFolder = value
			case "X-Origin":
				email.XOrigin = value
			case "X-FileName":
				email.XFileName = value
			}
		} else {
			email.Body += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return Email{}, err
	}

	return email, nil
}

func sendHTTPRequest(data []byte) error {
	url := "https://api.openobserve.ai/api/yair_camilo_organization_15016_FGL9xXNbnwFSWC2/default/_json"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creando la petición HTTP: %v", err)
	}
	req.SetBasicAuth("yaircamilo05@gmail.com", "2e0Z549386OPoci1n7BR")
	req.Header.Set("Content-Type", "application/json")

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
