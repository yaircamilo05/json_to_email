//cambiar el json por ndjson
//no guardar el archivo sino directamente en el codigo
//el tamaño del ndjson no puede ser mayor a 10mb
// hacerlo por curl

package main

import (
	"bufio"
	"encoding/json"
	//"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"sync"
	"time"
)

type Email struct {
	MessageID               string `json:"Message-ID"`
	Date                    string `json:"Date_"`
	From                    string `json:"From_"`
	To                      string `json:"To_"`
	Subject                 string `json:"Subject_"`
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
	n := time.Now()
	processDirectoriesRoot(dirPath)
	fmt.Println(time.Since(n))
}

func processDirectoriesRoot(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error al leer el directorio:", err)
		return
	}
	//encontramos las carpetas
	json := []string{}
	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		BuscarCarpetas(fullPath, &json)
	}
	fmt.Println(len(json))

	//hacemos un ndjson por cada carpeta
	var wg sync.WaitGroup
	for _, entry := range json {
		fmt.Println(entry)
	 	wg.Add(1)
	 	go func(entry string) {
	 		defer wg.Done()
	 		BuscarArchivos(entry)
	 	}(entry)
	 }

	 wg.Wait()
}

func BuscarCarpetas(path string, json *[]string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		subPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			if err := BuscarCarpetas(subPath, json); err != nil {
				return err
			}
			*json = append(*json, subPath)
		}
	}
	return nil
}


func BuscarArchivos(path string) {
	uploda_data := []Email{}
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error al leer el archivo", err)
		return
	}


	for _, entry := range entries {
		if !entry.IsDir() {
			subPath := filepath.Join(path, entry.Name())
			convertidorNdjson(subPath, &uploda_data)
			dataJson, err := json.Marshal(uploda_data)
			if err != nil {
				fmt.Println("Error al convertir los datos a JSON:", err)
				return
			}
			peticionCurl(string(dataJson))
		}
	}
}

func peticionCurl(data string) {

	cmd := exec.Command("curl", "-u", "yaircamilo05@gmail.com:6B4OG0158f27SZ3Av9Cr", "-k", "https://api.openobserve.ai/api/yair_camilo_organization_15016_FGL9xXNbnwFSWC2/emailsA/_bulk", "-d", "@"+data)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error ejecutando el comando para el archivo"+data+":", err)
		return
	}
	fmt.Println("Comando ejecutado exitosamente para el archivo" + data)
	fmt.Printf("Salida del comando:/n%s/n", output)
}

func convertidorNdjson(dataPath string, uploda_data *[]Email) {
	file, err := os.Open(dataPath)
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
			email.Body += line + "/n"
		}
		}
		*uploda_data = append(*uploda_data, email)
		return
	}
