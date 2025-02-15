package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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
	start := time.Now() // Inicia el temporizador

	dirPath := filepath.Join("..", "enron_mail_20110402", "maildir")
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

	elapsed := time.Since(start) // Calcula el tiempo transcurrido
	fmt.Printf("El proceso completo tomó %s\n", elapsed)
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
			processDirectory(dir)
		}(dir)
	}
	wg.Wait()
}

func processDirectory(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error al leer el directorio:", err)
		return
	}

	var emails []Email
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(path, entry.Name())
			email, err := readEmail(filePath)
			if err != nil {
				fmt.Println("Error al leer el archivo:", err)
				continue
			}
			emails = append(emails, email)
		}
	}

	err = sendEmailsAsJSON(emails)
	if err != nil {
		fmt.Println("Error al enviar el correo:", err)
	}
	time.Sleep(1000 * time.Millisecond) // Agregar un retraso de 1000 ms entre solicitudes
}

func readEmail(filePath string) (Email, error) {
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

func sendEmailsAsJSON(emails []Email) error {
	data := map[string]interface{}{
		"index":   "prueba",
		"records": emails,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error al convertir los correos a JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulkv2", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("error creando la petición HTTP: %v", err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
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

	fmt.Printf("Correo enviado exitosamente: %s\n", resp.Status)
	return nil
}
