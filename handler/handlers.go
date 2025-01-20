package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yaircamilo05/email_to_json/services"
)

type ProcessRequest struct {
	DirPath    string `json:"dir_path"`
	StreamName string `json:"stream_name"`
}

func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	var req ProcessRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.ProcessEmails(req.DirPath, req.StreamName)
	if err != nil {
		http.Error(w, "Error al procesar los emails: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Procesamiento completado para el stream: " + req.StreamName))
}

func GetEmailsHandler(w http.ResponseWriter, r *http.Request) {
	emails, err := services.GetEmails()
	if err != nil {
		http.Error(w, "Error al obtener los emails: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(emails)
}
