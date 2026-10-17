package logic

import (
	"errors"
	"strings"
)

var (
	ErrTituloVacio          = errors.New("el título del juego no puede estar vacío")
	ErrGeneroVacio          = errors.New("el género del juego no puede estar vacío")
	ErrLogrosNegativos      = errors.New("los logros totales no pueden ser negativos")
	ErrPlataformaVacia      = errors.New("la plataforma no puede estar vacía")
	ErrHorasNegativas       = errors.New("las horas jugadas no pueden ser negativas")
	ErrLogrosObtNegativos   = errors.New("los logros obtenidos no pueden ser negativos")
	ErrLogrosInconsistentes = errors.New("los logros obtenidos no pueden superar los logros totales")
)

type JuegoInput struct {
	Titulo        string `json:"titulo"`
	Genero        string `json:"genero"`
	LogrosTotales int32  `json:"logros_totales"`
}

func ValidateJuego(j JuegoInput) error {
	if strings.TrimSpace(j.Titulo) == "" {
		return ErrTituloVacio
	}
	if strings.TrimSpace(j.Genero) == "" {
		return ErrGeneroVacio
	}
	if j.LogrosTotales < 0 {
		return ErrLogrosNegativos
	}
	return nil
}

type CopiaInput struct {
	JuegoID         int32  `json:"juego_id"`
	Plataforma      string `json:"plataforma"`
	HorasJugadas    string `json:"horas_jugadas"`
	LogrosObtenidos int32  `json:"logros_obtenidos"`
}

func ValidateCopia(c CopiaInput) error {
	if c.JuegoID <= 0 {
		return errors.New("el ID del juego asociado debe ser válido")
	}
	if strings.TrimSpace(c.Plataforma) == "" {
		return ErrPlataformaVacia
	}
	if c.LogrosObtenidos < 0 {
		return ErrLogrosObtNegativos
	}
	return nil
}
