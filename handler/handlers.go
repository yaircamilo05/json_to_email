package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yaircamilo05/email_to_json/models"
	"github.com/yaircamilo05/email_to_json/services"
)

func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ProcessRequest
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
	var req models.GetAllEmailsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := services.GetEmails(req)
	if err != nil {
		http.Error(w, "Error al obtener los emails: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error al convertir a JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
