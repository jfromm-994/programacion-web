package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"programacion-web/logic"
)

type CopiasHandler struct {
	service *logic.Service
}

func NewCopiasHandler(s *logic.Service) *CopiasHandler {
	return &CopiasHandler{service: s}
}

func (h *CopiasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) == 1 && parts[0] == "copias" {
		switch r.Method {
		case http.MethodGet:
			h.listCopiasByJuego(w, r)
		case http.MethodPost:
			h.createCopia(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 2 && parts[0] == "copias" {
		parsedID, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil || parsedID <= 0 {
			http.Error(w, "ID de copia inválido", http.StatusBadRequest)
			return
		}
		id := int32(parsedID)

		if r.Method == http.MethodDelete {
			h.deleteCopia(w, r, id)
			return
		}
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

func (h *CopiasHandler) listCopiasByJuego(w http.ResponseWriter, r *http.Request) {
	juegoIDStr := r.URL.Query().Get("juego_id")
	if juegoIDStr == "" {
		http.Error(w, "Falta el parámetro de consulta ?juego_id=", http.StatusBadRequest)
		return
	}

	parsedJuegoID, err := strconv.ParseInt(juegoIDStr, 10, 32)
	if err != nil || parsedJuegoID <= 0 {
		http.Error(w, "Parámetro juego_id inválido", http.StatusBadRequest)
		return
	}

	copias, err := h.service.ListCopiasByJuego(r.Context(), int32(parsedJuegoID))
	if err != nil {
		http.Error(w, "Error al listar copias", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(copias)
}

func (h *CopiasHandler) createCopia(w http.ResponseWriter, r *http.Request) {
	var input logic.CopiaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo JSON inválido", http.StatusBadRequest)
		return
	}

	copia, err := h.service.CreateCopia(r.Context(), input)
	if err != nil {
		if errors.Is(err, logic.ErrJuegoIDInvalido) || errors.Is(err, logic.ErrPlataformaVacia) || errors.Is(err, logic.ErrLogrosObtNegativos) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Error al guardar la copia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(copia)
}

func (h *CopiasHandler) deleteCopia(w http.ResponseWriter, r *http.Request, id int32) {
	if err := h.service.DeleteCopia(r.Context(), id); err != nil {
		http.Error(w, "Error al eliminar la copia", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
