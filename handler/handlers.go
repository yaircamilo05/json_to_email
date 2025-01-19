package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yaircamilo05/email_to_json/utils"
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

	directories := []string{}
	err = utils.FindDirectories(req.DirPath, &directories)
	if err != nil {
		http.Error(w, "Error al buscar carpetas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	directoriesWithFiles := utils.FilterDirectoriesWithFiles(directories)
	utils.Core(directoriesWithFiles, req.StreamName)

	w.Write([]byte("Procesamiento completado para el stream: " + req.StreamName))
}
