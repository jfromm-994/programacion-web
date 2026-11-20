package main

import (
	"context"
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

	// Precarga de registros iniciales si la BD está vacía
	if err := service.SeedInitialData(context.Background()); err != nil {
		log.Printf("Advertencia al precargar datos iniciales: %v", err)
	} else {
		log.Println("Datos iniciales verificados/precargados correctamente.")
	}

	// Controladores específicos
	juegosCtrl := handlers.NewJuegosController(service)
	copiasCtrl := handlers.NewCopiasController(service)

	mux := http.NewServeMux()

	// Rutas de Catálogo y Juegos
	mux.HandleFunc("GET /", juegosCtrl.CatalogoHandler)
	mux.HandleFunc("GET /juegos", juegosCtrl.CatalogoHandler)
	mux.HandleFunc("POST /juegos", juegosCtrl.CatalogoHandler)
	mux.HandleFunc("GET /juegos/detalle", juegosCtrl.DetalleHandler)
	mux.HandleFunc("DELETE /juegos/{id}", juegosCtrl.DeleteHandler)
	mux.HandleFunc("POST /juegos/delete", juegosCtrl.DeleteHandler) // Fallback sin JS

	// Rutas de Copias
	mux.HandleFunc("POST /copias", copiasCtrl.CreateHandler)
	mux.HandleFunc("DELETE /copias/{id}", copiasCtrl.DeleteHandler)
	mux.HandleFunc("POST /copias/delete", copiasCtrl.DeleteHandler) // Fallback sin JS

	// Middleware de observabilidad/logging
	loggedMux := handlers.LoggingMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor GameVault (TP6 HTMX) corriendo en http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, loggedMux); err != nil {
		log.Fatalf("Error en el servidor HTTP: %v", err)
	}
}
