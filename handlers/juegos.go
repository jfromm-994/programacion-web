package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"programacion-web/logic"
	"programacion-web/views"
)

type JuegosHandler struct {
	service *logic.Service
}

func NewJuegosHandler(s *logic.Service) *JuegosHandler {
	return &JuegosHandler{service: s}
}

func (h *JuegosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// Ruta principal: "/" o "/juegos"
	if path == "" || (len(parts) == 1 && parts[0] == "juegos") {
		switch r.Method {
		case http.MethodGet:
			h.renderCatalogo(w, r, "")
		case http.MethodPost:
			h.createJuegoPRG(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Detalle de un juego y sus copias: "/juegos/detalle"
	if len(parts) == 2 && parts[0] == "juegos" && parts[1] == "detalle" {
		if r.Method == http.MethodGet {
			h.renderDetalle(w, r, "")
			return
		}
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Eliminación de un juego: "/juegos/delete" (vía formulario POST sin JS)
	if len(parts) == 2 && parts[0] == "juegos" && parts[1] == "delete" {
		if r.Method == http.MethodPost {
			h.deleteJuegoPRG(w, r)
			return
		}
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

func (h *JuegosHandler) renderCatalogo(w http.ResponseWriter, r *http.Request, errorMsg string) {
	juegos, err := h.service.ListJuegosConMetricas(r.Context())
	if err != nil {
		http.Error(w, "Error al recuperar juegos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = views.JuegosPage(juegos, errorMsg).Render(r.Context(), w)
}

func (h *JuegosHandler) renderDetalle(w http.ResponseWriter, r *http.Request, errorMsg string) {
	idStr := r.URL.Query().Get("id")
	parsedID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || parsedID <= 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	juegoID := int32(parsedID)
	juego, err := h.service.GetJuegoByID(r.Context(), juegoID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	copias, err := h.service.ListCopiasByJuego(r.Context(), juegoID)
	if err != nil {
		http.Error(w, "Error al obtener copias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = views.CopiasPage(juego, copias, errorMsg).Render(r.Context(), w)
}

func (h *JuegosHandler) createJuegoPRG(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderCatalogo(w, r, "Formulario inválido")
		return
	}

	titulo := strings.TrimSpace(r.FormValue("titulo"))
	genero := strings.TrimSpace(r.FormValue("genero"))
	logrosStr := strings.TrimSpace(r.FormValue("logros_totales"))

	var logros int32
	if logrosStr != "" {
		val, err := strconv.ParseInt(logrosStr, 10, 32)
		if err != nil {
			h.renderCatalogo(w, r, "Los logros deben ser un número válido")
			return
		}
		logros = int32(val)
	}

	input := logic.JuegoInput{
		Titulo:        titulo,
		Genero:        genero,
		LogrosTotales: logros,
	}

	_, err := h.service.CreateJuego(r.Context(), input)
	if err != nil {
		h.renderCatalogo(w, r, err.Error())
		return
	}

	// Patrón Post/Redirect/Get
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *JuegosHandler) deleteJuegoPRG(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	parsedID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || parsedID <= 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_ = h.service.DeleteJuego(r.Context(), int32(parsedID))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
