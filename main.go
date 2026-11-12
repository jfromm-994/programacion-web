package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "programacion-web/db/sqlc"
	"programacion-web/handlers"
	"programacion-web/logic"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/gamevault?sslmode=disable"
	}

	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Error al abrir conexión con la base de datos: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}
	log.Println("Conexión con PostgreSQL establecida exitosamente.")

	queries := db.New(conn)
	service := logic.NewService(queries)

	// Controladores con vistas SSR y lógica de negocio
	juegosHandler := handlers.NewJuegosHandler(service)
	copiasHandler := handlers.NewCopiasHandler(service)

	mux := http.NewServeMux()

	// Rutas de la aplicación web (SSR con Templ)
	mux.Handle("/", juegosHandler)
	mux.Handle("/juegos", juegosHandler)
	mux.Handle("/juegos/", juegosHandler)
	mux.Handle("/copias", copiasHandler)
	mux.Handle("/copias/", copiasHandler)

	// Middleware de observabilidad/logging requerido por la cátedra
	loggedMux := handlers.LoggingMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor GameVault corriendo en http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, loggedMux); err != nil {
		log.Fatalf("Error en el servidor HTTP: %v", err)
	}
}
