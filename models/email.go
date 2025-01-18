package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func ProcessDirectory(path string) {
	emails := []Email{}
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error al leer el directorio:", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(path, entry.Name())
			email, err := ReadEmail(filePath)
			if err != nil {
				fmt.Println("Error al leer el archivo:", err)
				continue
			}
			emails = append(emails, email)
		}
	}

	jsonData, err := json.Marshal(emails)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	err = SendHTTPRequest(jsonData)
	if err != nil {
		fmt.Println("Error al enviar la petición HTTP:", err)
		if strings.Contains(err.Error(), "413 Request Entity Too Large") {
			HandleLargeRequest(emails)
		}
	}
}

func ReadEmail(filePath string) (Email, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Email{}, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	email := Email{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return Email{}, err
		}
		line = strings.TrimSpace(line)
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

	return email, nil
}

func SendHTTPRequest(data []byte) error {
	url := "http://localhost:5080/api/prueba/emails/_json"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creando la petición HTTP: %v", err)
	}
	req.SetBasicAuth("root@example.com", "T9uplVBu16xjKUrd")
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

func HandleLargeRequest(emails []Email) {
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

	err = SendHTTPRequest(jsonData1)
	if err != nil {
		fmt.Println("Error al enviar la petición HTTP:", err)
	}
	fmt.Println("Petición HTTP enviada exitosamente en el segundo intento")

	err = SendHTTPRequest(jsonData2)
	if err != nil {
		fmt.Println("Error al enviar la petición HTTP:", err)
	}
	fmt.Println("Petición HTTP enviada exitosamente en el segundo intento")
}
