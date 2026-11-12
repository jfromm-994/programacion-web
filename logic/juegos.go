package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	db "programacion-web/db/sqlc"
)

var (
	ErrNotFound           = errors.New("recurso no encontrado")
	ErrTituloVacio        = errors.New("el título del juego no puede estar vacío")
	ErrGeneroVacio        = errors.New("el género del juego no puede estar vacío")
	ErrLogrosNegativos    = errors.New("los logros totales no pueden ser negativos")
	ErrPlataformaVacia    = errors.New("la plataforma no puede estar vacía")
	ErrLogrosObtNegativos = errors.New("los logros obtenidos no pueden ser negativos")
	ErrJuegoIDInvalido    = errors.New("el ID del juego asociado debe ser válido")
)

type JuegoInput struct {
	Titulo        string `json:"titulo"`
	Genero        string `json:"genero"`
	LogrosTotales int32  `json:"logros_totales"`
}

type CopiaInput struct {
	JuegoID         int32  `json:"juego_id"`
	Plataforma      string `json:"plataforma"`
	HorasJugadas    string `json:"horas_jugadas"`
	LogrosObtenidos int32  `json:"logros_obtenidos"`
}

type Service struct {
	queries *db.Queries
}

func NewService(q *db.Queries) *Service {
	return &Service{queries: q}
}

// =================== Lógica de Juegos ===================

func (s *Service) ListJuegos(ctx context.Context) ([]db.Juego, error) {
	return s.queries.ListJuegos(ctx)
}

func (s *Service) GetJuegoByID(ctx context.Context, id int32) (db.Juego, error) {
	juego, err := s.queries.GetJuegoByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Juego{}, ErrNotFound
		}
		return db.Juego{}, err
	}
	return juego, nil
}

func (s *Service) CreateJuego(ctx context.Context, input JuegoInput) (db.Juego, error) {
	if strings.TrimSpace(input.Titulo) == "" {
		return db.Juego{}, ErrTituloVacio
	}
	if strings.TrimSpace(input.Genero) == "" {
		return db.Juego{}, ErrGeneroVacio
	}
	if input.LogrosTotales < 0 {
		return db.Juego{}, ErrLogrosNegativos
	}

	return s.queries.CreateJuego(ctx, db.CreateJuegoParams{
		Titulo:        input.Titulo,
		Genero:        input.Genero,
		LogrosTotales: input.LogrosTotales,
	})
}

func (s *Service) UpdateJuego(ctx context.Context, id int32, input JuegoInput) (db.Juego, error) {
	_, err := s.GetJuegoByID(ctx, id)
	if err != nil {
		return db.Juego{}, err
	}

	if strings.TrimSpace(input.Titulo) == "" {
		return db.Juego{}, ErrTituloVacio
	}
	if strings.TrimSpace(input.Genero) == "" {
		return db.Juego{}, ErrGeneroVacio
	}
	if input.LogrosTotales < 0 {
		return db.Juego{}, ErrLogrosNegativos
	}

	err = s.queries.UpdateJuego(ctx, db.UpdateJuegoParams{
		ID:            id,
		Titulo:        input.Titulo,
		Genero:        input.Genero,
		LogrosTotales: input.LogrosTotales,
	})
	if err != nil {
		return db.Juego{}, err
	}

	return s.queries.GetJuegoByID(ctx, id)
}

func (s *Service) DeleteJuego(ctx context.Context, id int32) error {
	_, err := s.GetJuegoByID(ctx, id)
	if err != nil {
		return err
	}
	return s.queries.DeleteJuego(ctx, id)
}

// =================== Lógica de Copias ===================

func (s *Service) ListCopiasByJuego(ctx context.Context, juegoID int32) ([]db.Copia, error) {
	return s.queries.ListCopiasByJuego(ctx, juegoID)
}

func (s *Service) CreateCopia(ctx context.Context, input CopiaInput) (db.Copia, error) {
	if input.JuegoID <= 0 {
		return db.Copia{}, ErrJuegoIDInvalido
	}
	if strings.TrimSpace(input.Plataforma) == "" {
		return db.Copia{}, ErrPlataformaVacia
	}
	if input.LogrosObtenidos < 0 {
		return db.Copia{}, ErrLogrosObtNegativos
	}

	return s.queries.CreateCopia(ctx, db.CreateCopiaParams{
		JuegoID:         input.JuegoID,
		Plataforma:      input.Plataforma,
		HorasJugadas:    input.HorasJugadas,
		LogrosObtenidos: input.LogrosObtenidos,
	})
}

func (s *Service) DeleteCopia(ctx context.Context, id int32) error {
	return s.queries.DeleteCopia(ctx, id)
}
