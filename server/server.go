package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"

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
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutos
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/process", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Create("cpu_profile.prof")
		if err != nil {
			http.Error(w, "No se pudo crear el archivo de perfil de CPU", http.StatusInternalServerError)
			return
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			http.Error(w, "No se pudo iniciar el perfil de CPU", http.StatusInternalServerError)
			return
		}

		done := make(chan struct{})
		go func() {
			handler.ProcessHandler(w, r)
			close(done)
		}()

		go func() {
			time.Sleep(1 * time.Minute)
			pprof.StopCPUProfile()
			f.Close()
			fmt.Println("Perfil de CPU generado en cpu_profile.prof")

			cmd := exec.Command("go", "tool", "pprof", "-png", "cpu_profile.prof")
			output, err := cmd.Output()
			if err != nil {
				log.Fatalf("Error generando el gráfico del perfil de CPU: %v", err)
			}
			err = os.WriteFile("cpu_profile.png", output, 0644)
			if err != nil {
				log.Fatalf("Error guardando el gráfico del perfil de CPU: %v", err)
			}
			fmt.Println("Gráfico del perfil de CPU generado en cpu_profile.png")
		}()

		<-done
	})

	// Consigue todos los emails de la base de datos
	r.Post("/get_emails", handler.GetEmailsHandler)

	fmt.Println("Servidor escuchando en el puerto 3000")
	http.ListenAndServe(":3000", r)
}
