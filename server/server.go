package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	// Configurar CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"}, // Reemplaza con la URL de tu frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutos
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/process", handler.ProcessHandler)
	// Consigue todos los emails de la base de datos
	r.Post("/get_emails", handler.GetEmailsHandler)

	fmt.Println("Servidor escuchando en el puerto 3000")
	http.ListenAndServe(":3000", r)
}
