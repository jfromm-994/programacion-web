package main

import (
	"context"
	"database/sql"
	"testing"

	db "programacion-web/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestGameVaultPersistence(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost:5432/gamevault?sslmode=disable"
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("No se pudo configurar la conexión: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		t.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}

	queries := db.New(conn)
	ctx := context.Background()

	// 1. Prueba: Crear Juego
	juego, err := queries.CreateJuego(ctx, db.CreateJuegoParams{
		Titulo:        "Celeste",
		Genero:        "Platformer",
		LogrosTotales: 30,
	})
	if err != nil {
		t.Fatalf("Error al crear juego: %v", err)
	}
	if juego.ID == 0 {
		t.Errorf("Se esperaba un ID válido mayor a 0, se obtuvo %d", juego.ID)
	}

	// 2. Prueba: Crear Copia asociada
	copia, err := queries.CreateCopia(ctx, db.CreateCopiaParams{
		JuegoID:         juego.ID,
		Plataforma:      "Nintendo Switch",
		HorasJugadas:    "25.00",
		LogrosObtenidos: 20,
	})
	if err != nil {
		t.Fatalf("Error al crear copia: %v", err)
	}
	if copia.ID == 0 {
		t.Errorf("Se esperaba un ID válido de copia, se obtuvo %d", copia.ID)
	}

	// 3. Prueba: Listar Juegos
	juegos, err := queries.ListJuegos(ctx)
	if err != nil {
		t.Fatalf("Error al listar juegos: %v", err)
	}
	if len(juegos) == 0 {
		t.Errorf("Se esperaba al menos 1 juego en la lista")
	}

	// 4. Prueba: Borrar Juego (Eliminación en Cascada)
	err = queries.DeleteJuego(ctx, juego.ID)
	if err != nil {
		t.Fatalf("Error al eliminar juego: %v", err)
	}
}
