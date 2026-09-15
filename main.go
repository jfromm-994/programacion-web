package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	db "programacion-web/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// 1. Cadena de conexión a PostgreSQL (coincide con docker-compose.yml)
	connStr := "postgres://postgres:postgres@localhost:5432/gamevault?sslmode=disable"

	// 2. Abrir conexión con la base de datos
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Error al configurar la conexión: %v", err)
	}
	defer conn.Close()

	// 3. Verificar que la base de datos responda
	if err := conn.Ping(); err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	fmt.Println("¡Conexión exitosa a PostgreSQL!")

	// 4. Instanciar el paquete generado por sqlc
	queries := db.New(conn)
	ctx := context.Background()

	// 5. Probar crear un Juego de prueba
	nuevoJuego, err := queries.CreateJuego(ctx, db.CreateJuegoParams{
		Titulo:        "Hollow Knight",
		Genero:        "Metroidvania",
		LogrosTotales: 63,
	})
	if err != nil {
		log.Printf("Aviso/Error al insertar juego: %v", err)
	} else {
		fmt.Printf("Juego CREADO con éxito -> ID: %d | Título: %s\n", nuevoJuego.ID, nuevoJuego.Titulo)

		// 6. Probar crear una Copia para este juego
		nuevaCopia, err := queries.CreateCopia(ctx, db.CreateCopiaParams{
			JuegoID:         nuevoJuego.ID,
			Plataforma:      "Steam",
			HorasJugadas:    "45.50",
			LogrosObtenidos: 50,
		})
		if err != nil {
			log.Printf("Error al insertar copia: %v", err)
		} else {
			fmt.Printf("Copia CREADA con éxito -> ID: %d | Plataforma: %s\n", nuevaCopia.ID, nuevaCopia.Plataforma)
		}
	}

	// 7. Listar todos los juegos almacenados
	juegos, err := queries.ListJuegos(ctx)
	if err != nil {
		log.Fatalf("Error al listar juegos: %v", err)
	}

	fmt.Println("\n--- Lista de Juegos en Game Vault ---")
	for _, j := range juegos {
		fmt.Printf("- [%d] %s (%s) - Logros: %d\n", j.ID, j.Titulo, j.Genero, j.LogrosTotales)
	}
}
