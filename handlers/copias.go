package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"programacion-web/logic"
	"programacion-web/views"
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

	// Alta de copia: "/copias" vía POST
	if len(parts) == 1 && parts[0] == "copias" {
		if r.Method == http.MethodPost {
			h.createCopiaPRG(w, r)
			return
		}
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Eliminación de copia: "/copias/delete" vía POST
	if len(parts) == 2 && parts[0] == "copias" && parts[1] == "delete" {
		if r.Method == http.MethodPost {
			h.deleteCopiaPRG(w, r)
			return
		}
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

func (h *CopiasHandler) createCopiaPRG(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	juegoIDStr := r.FormValue("juego_id")
	plataforma := strings.TrimSpace(r.FormValue("plataforma"))
	horasStr := strings.TrimSpace(r.FormValue("horas_jugadas"))
	logrosStr := strings.TrimSpace(r.FormValue("logros_obtenidos"))

	parsedJuegoID, err := strconv.ParseInt(juegoIDStr, 10, 32)
	if err != nil || parsedJuegoID <= 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	juegoID := int32(parsedJuegoID)

	var logros int32
	if logrosStr != "" {
		val, err := strconv.ParseInt(logrosStr, 10, 32)
		if err == nil {
			logros = int32(val)
		}
	}

	if horasStr == "" {
		horasStr = "0.0"
	}

	input := logic.CopiaInput{
		JuegoID:         juegoID,
		Plataforma:      plataforma,
		HorasJugadas:    horasStr,
		LogrosObtenidos: logros,
	}

	_, err = h.service.CreateCopia(r.Context(), input)
	if err != nil {
		// En caso de error, re-renderizamos la vista de copias mostrando el cartel de error
		juego, getErr := h.service.GetJuegoByID(r.Context(), juegoID)
		if getErr != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		copias, _ := h.service.ListCopiasByJuego(r.Context(), juegoID)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = views.CopiasPage(juego, copias, err.Error()).Render(r.Context(), w)
		return
	}

	// Éxito: Post/Redirect/Get hacia la vista de detalle
	redirectURL := fmt.Sprintf("/juegos/detalle?id=%d", juegoID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *CopiasHandler) deleteCopiaPRG(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	juegoIDStr := r.FormValue("juego_id")

	parsedID, err := strconv.ParseInt(idStr, 10, 32)
	if err == nil && parsedID > 0 {
		_ = h.service.DeleteCopia(r.Context(), int32(parsedID))
	}

	if juegoIDStr != "" {
		http.Redirect(w, r, fmt.Sprintf("/juegos/detalle?id=%s", juegoIDStr), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
