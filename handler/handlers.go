package handler

import (
	"net/http"

	"github.com/yaircamilo05/email_to_json/utils"
)

func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	dirPath := "../enron_mail_20110402/maildir"
	directories := []string{}
	err := utils.FindDirectories(dirPath, &directories)
	if err != nil {
		http.Error(w, "Error al buscar carpetas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	directoriesWithFiles := utils.FilterDirectoriesWithFiles(directories)
	utils.Core(directoriesWithFiles)

	w.Write([]byte("Procesamiento completado"))
}
