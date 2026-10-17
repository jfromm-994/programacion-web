package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	db "programacion-web/db/sqlc"
	"programacion-web/logic"
)

type JuegosHandler struct {
	queries *db.Queries
}

func NewJuegosHandler(q *db.Queries) *JuegosHandler {
	return &JuegosHandler{queries: q}
}

func (h *JuegosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// Ruta base: /juegos
	if len(parts) == 1 && parts[0] == "juegos" {
		switch r.Method {
		case http.MethodGet:
			h.listJuegos(w, r)
		case http.MethodPost:
			h.createJuego(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Ruta con ID: /juegos/{id}
	if len(parts) == 2 && parts[0] == "juegos" {
		parsedID, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil || parsedID <= 0 {
			http.Error(w, "ID de juego inválido", http.StatusBadRequest)
			return
		}
		id := int32(parsedID)

		switch r.Method {
		case http.MethodGet:
			h.getJuego(w, r, id)
		case http.MethodPut:
			h.updateJuego(w, r, id)
		case http.MethodDelete:
			h.deleteJuego(w, r, id)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func (h *JuegosHandler) listJuegos(w http.ResponseWriter, r *http.Request) {
	juegos, err := h.queries.ListJuegos(r.Context())
	if err != nil {
		http.Error(w, "Error al listar juegos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(juegos)
}

func (h *JuegosHandler) getJuego(w http.ResponseWriter, r *http.Request, id int32) {
	juego, err := h.queries.GetJuegoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Juego no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener juego", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(juego)
}

func (h *JuegosHandler) createJuego(w http.ResponseWriter, r *http.Request) {
	var input logic.JuegoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo JSON inválido", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateJuego(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	juego, err := h.queries.CreateJuego(r.Context(), db.CreateJuegoParams{
		Titulo:        input.Titulo,
		Genero:        input.Genero,
		LogrosTotales: input.LogrosTotales,
	})
	if err != nil {
		http.Error(w, "Error al guardar el juego en base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(juego)
}

func (h *JuegosHandler) updateJuego(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetJuegoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Juego no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	var input logic.JuegoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo JSON inválido", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateJuego(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.queries.UpdateJuego(r.Context(), db.UpdateJuegoParams{
		ID:            id,
		Titulo:        input.Titulo,
		Genero:        input.Genero,
		LogrosTotales: input.LogrosTotales,
	})
	if err != nil {
		http.Error(w, "Error al actualizar juego", http.StatusInternalServerError)
		return
	}

	juegoActualizado, _ := h.queries.GetJuegoByID(r.Context(), id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(juegoActualizado)
}

func (h *JuegosHandler) deleteJuego(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetJuegoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Juego no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	if err := h.queries.DeleteJuego(r.Context(), id); err != nil {
		http.Error(w, "Error al eliminar juego", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
