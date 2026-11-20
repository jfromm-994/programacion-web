package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	db "programacion-web/db/sqlc"
)

var (
	ErrNotFound            = errors.New("recurso no encontrado")
	ErrTituloVacio         = errors.New("el título del juego no puede estar vacío")
	ErrGeneroVacio         = errors.New("el género del juego no puede estar vacío")
	ErrLogrosNegativos     = errors.New("los logros totales no pueden ser negativos")
	ErrPlataformaVacia     = errors.New("la plataforma no puede estar vacía")
	ErrLogrosObtNegativos  = errors.New("los logros obtenidos no pueden ser negativos")
	ErrJuegoIDInvalido     = errors.New("el ID del juego asociado debe ser válido")
	ErrLogrosExcedidos     = errors.New("los logros obtenidos no pueden superar los logros totales del juego")
	ErrPlataformaDuplicada = errors.New("ya tienes registrada una copia para esa plataforma")
)

// cleanText normaliza cadenas de texto en toda la capa lógica (recorta espacios y aplica formato título)
func cleanText(s string) string {
	return cases.Title(language.Und).String(strings.ToLower(strings.TrimSpace(s)))
}

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

func (s *Service) ListJuegosConMetricas(ctx context.Context) ([]db.ListJuegosConMetricasRow, error) {
	return s.queries.ListJuegosConMetricas(ctx)
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

func (s *Service) GetJuegoByTitulo(ctx context.Context, titulo string) (db.Juego, error) {
	tituloLimpio := cleanText(titulo)
	juego, err := s.queries.GetJuegoByTitulo(ctx, tituloLimpio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Juego{}, ErrNotFound
		}
		return db.Juego{}, err
	}
	return juego, nil
}

func (s *Service) CreateJuego(ctx context.Context, input JuegoInput) (db.Juego, error) {
	titulo := cleanText(input.Titulo)
	if titulo == "" {
		return db.Juego{}, ErrTituloVacio
	}
	genero := cleanText(input.Genero)
	if genero == "" {
		return db.Juego{}, ErrGeneroVacio
	}
	if input.LogrosTotales < 0 {
		return db.Juego{}, ErrLogrosNegativos
	}

	juego, err := s.queries.CreateJuego(ctx, db.CreateJuegoParams{
		Titulo:        titulo,
		Genero:        genero,
		LogrosTotales: input.LogrosTotales,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			return db.Juego{}, errors.New("ya existe un juego registrado con ese título")
		}
		return db.Juego{}, err
	}
	return juego, nil
}

func (s *Service) UpdateJuego(ctx context.Context, id int32, input JuegoInput) (db.Juego, error) {
	_, err := s.GetJuegoByID(ctx, id)
	if err != nil {
		return db.Juego{}, err
	}

	titulo := cleanText(input.Titulo)
	if titulo == "" {
		return db.Juego{}, ErrTituloVacio
	}
	genero := cleanText(input.Genero)
	if genero == "" {
		return db.Juego{}, ErrGeneroVacio
	}
	if input.LogrosTotales < 0 {
		return db.Juego{}, ErrLogrosNegativos
	}

	err = s.queries.UpdateJuego(ctx, db.UpdateJuegoParams{
		ID:            id,
		Titulo:        titulo,
		Genero:        genero,
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
	plataforma := cleanText(input.Plataforma)
	if plataforma == "" {
		return db.Copia{}, ErrPlataformaVacia
	}
	if input.LogrosObtenidos < 0 {
		return db.Copia{}, ErrLogrosObtNegativos
	}

	// Validación de dominio: consultar juego base para verificar tope de logros
	juego, err := s.GetJuegoByID(ctx, input.JuegoID)
	if err != nil {
		return db.Copia{}, err
	}

	if input.LogrosObtenidos > juego.LogrosTotales {
		return db.Copia{}, ErrLogrosExcedidos
	}

	copia, err := s.queries.CreateCopia(ctx, db.CreateCopiaParams{
		JuegoID:         input.JuegoID,
		Plataforma:      plataforma,
		HorasJugadas:    input.HorasJugadas,
		LogrosObtenidos: input.LogrosObtenidos,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			return db.Copia{}, ErrPlataformaDuplicada
		}
		return db.Copia{}, err
	}

	return copia, nil
}

func (s *Service) DeleteCopia(ctx context.Context, id int32) error {
	return s.queries.DeleteCopia(ctx, id)
}

// SeedInitialData carga registros de prueba si la base de datos está vacía (recomendación de cátedra)
func (s *Service) SeedInitialData(ctx context.Context) error {
	juegos, err := s.ListJuegos(ctx)
	if err != nil {
		return err
	}

	if len(juegos) > 0 {
		return nil
	}

	// 1. Hollow Knight
	hk, err := s.CreateJuego(ctx, JuegoInput{
		Titulo:        "Hollow Knight",
		Genero:        "Metroidvania",
		LogrosTotales: 63,
	})
	if err == nil {
		_, _ = s.CreateCopia(ctx, CopiaInput{
			JuegoID:         hk.ID,
			Plataforma:      "PC Steam",
			HorasJugadas:    "45.50",
			LogrosObtenidos: 52,
		})
		_, _ = s.CreateCopia(ctx, CopiaInput{
			JuegoID:         hk.ID,
			Plataforma:      "Nintendo Switch",
			HorasJugadas:    "20.00",
			LogrosObtenidos: 30,
		})
	}

	// 2. Celeste
	celeste, err := s.CreateJuego(ctx, JuegoInput{
		Titulo:        "Celeste",
		Genero:        "Plataformas",
		LogrosTotales: 32,
	})
	if err == nil {
		_, _ = s.CreateCopia(ctx, CopiaInput{
			JuegoID:         celeste.ID,
			Plataforma:      "PC Steam",
			HorasJugadas:    "18.20",
			LogrosObtenidos: 25,
		})
	}

	// 3. Dead Cells
	dc, err := s.CreateJuego(ctx, JuegoInput{
		Titulo:        "Dead Cells",
		Genero:        "Roguelike",
		LogrosTotales: 105,
	})
	if err == nil {
		_, _ = s.CreateCopia(ctx, CopiaInput{
			JuegoID:         dc.ID,
			Plataforma:      "PlayStation 5",
			HorasJugadas:    "60.00",
			LogrosObtenidos: 70,
		})
	}

	return nil
}
