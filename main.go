package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	dirPath := "C:/Users/yairc/OneDrive/Documentos/universidad/semestre 7/truora 2/enron_mail_20110402/maildir"
	processDirectories(dirPath, "")
}

func processDirectories(path string, prefix string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error al leer el directorio:", err)
		return
	}

	var wg sync.WaitGroup

	for _, entry := range entries {
		subPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			go func(subPath string) {
				defer wg.Done()
				processDirectories(subPath, prefix+"  ")
			}(subPath)
		} else {
			go func(subPath string) {
				defer wg.Done()
				processFile(subPath)
			}(subPath)
		}
	}

	wg.Wait()
}

func processFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
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

	jsonData, err := json.MarshalIndent(email, "", "  ")
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	dirName := filepath.Base(filepath.Dir(filePath))

	// Crea un nuevo archivo en el directorio actual con el nombre del directorio padre
	jsonFile, err := os.Create("C:/Users/yairc/OneDrive/Documentos/universidad/semestre 7/truora 2/json/" + dirName + ".json")
	if err != nil {
		fmt.Println("Error al crear el archivo JSON:", err)
		return
	}
	defer jsonFile.Close()

	// Escribe los datos JSON en el archivo
	jsonFile.Write(jsonData)
	jsonFile.WriteString("\n")
}
