package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yaircamilo05/email_to_json/handler"
)

func main() {
	// Inicia el servidor HTTP para el profiling
	go func() {
		fmt.Println("Iniciando servidor de profiling en :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/process", handler.ProcessHandler)

	fmt.Println("Servidor escuchando en el puerto 3000")
	http.ListenAndServe(":3000", r)
}
