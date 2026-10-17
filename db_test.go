package main

import (
	"context"
	"database/sql"
	"testing"

	db "programacion-web/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// setupTestDB centraliza la conexión para todos los tests
func setupTestDB(t *testing.T) (*db.Queries, context.Context) {
	connStr := "postgres://postgres:postgres@localhost:5432/gamevault?sslmode=disable"
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("No se pudo configurar la conexión: %v", err)
	}

	if err := conn.Ping(); err != nil {
		t.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}

	// t.Cleanup asegura que la conexión se cierre al finalizar cada test
	t.Cleanup(func() {
		conn.Close()
	})

	return db.New(conn), context.Background()
}

func TestCreateJuego(t *testing.T) {
	queries, ctx := setupTestDB(t)

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

	// Limpieza posterior al test
	_ = queries.DeleteJuego(ctx, juego.ID)
}

func TestListJuegos(t *testing.T) {
	queries, ctx := setupTestDB(t)

	// Creamos un registro para asegurar que la lista no esté vacía
	juego, err := queries.CreateJuego(ctx, db.CreateJuegoParams{
		Titulo:        "Hollow Knight",
		Genero:        "Metroidvania",
		LogrosTotales: 63,
	})
	if err != nil {
		t.Fatalf("Error al preparar datos para ListJuegos: %v", err)
	}
	defer queries.DeleteJuego(ctx, juego.ID)

	juegos, err := queries.ListJuegos(ctx)
	if err != nil {
		t.Fatalf("Error al listar juegos: %v", err)
	}
	if len(juegos) == 0 {
		t.Errorf("Se esperaba al menos 1 juego en la lista")
	}
}

func TestCreateCopia(t *testing.T) {
	queries, ctx := setupTestDB(t)

	// Creamos el juego base necesario para asociar la clave foránea
	juego, err := queries.CreateJuego(ctx, db.CreateJuegoParams{
		Titulo:        "Dead Cells",
		Genero:        "Roguelike",
		LogrosTotales: 100,
	})
	if err != nil {
		t.Fatalf("Error al preparar juego para la copia: %v", err)
	}
	defer queries.DeleteJuego(ctx, juego.ID)

	copia, err := queries.CreateCopia(ctx, db.CreateCopiaParams{
		JuegoID:         juego.ID,
		Plataforma:      "Steam",
		HorasJugadas:    "15.50",
		LogrosObtenidos: 45,
	})
	if err != nil {
		t.Fatalf("Error al crear copia: %v", err)
	}
	if copia.ID == 0 {
		t.Errorf("Se esperaba un ID válido de copia, se obtuvo %d", copia.ID)
	}
}

func TestDeleteJuego(t *testing.T) {
	queries, ctx := setupTestDB(t)

	juego, err := queries.CreateJuego(ctx, db.CreateJuegoParams{
		Titulo:        "Hades",
		Genero:        "Action Roguelike",
		LogrosTotales: 49,
	})
	if err != nil {
		t.Fatalf("Error al crear juego para eliminar: %v", err)
	}

	err = queries.DeleteJuego(ctx, juego.ID)
	if err != nil {
		t.Fatalf("Error al eliminar juego: %v", err)
	}
}
