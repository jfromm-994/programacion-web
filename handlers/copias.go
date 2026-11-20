package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"programacion-web/logic"
	"programacion-web/views"
)

type CopiasController struct {
	service *logic.Service
}

func NewCopiasController(s *logic.Service) *CopiasController {
	return &CopiasController{service: s}
}

// CreateHandler: POST /copias
func (c *CopiasController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"

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

	_, err = c.service.CreateCopia(r.Context(), input)
	if err != nil {
		if isHTMX {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_ = views.CopiaFormOOB(juegoID, err.Error()).Render(r.Context(), w)
			return
		}
		redirectURL := fmt.Sprintf("/juegos/detalle?id=%d&error=%s", juegoID, url.QueryEscape(err.Error()))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Éxito: Si es HTMX, devolvemos la tabla de copias actualizada + formulario limpio OOB
	if isHTMX {
		juego, _ := c.service.GetJuegoByID(r.Context(), juegoID)
		copias, _ := c.service.ListCopiasByJuego(r.Context(), juegoID)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = views.CopiasTable(juego, copias).Render(r.Context(), w)
		_ = views.CopiaFormOOB(juegoID, "").Render(r.Context(), w)
		return
	}

	redirectURL := fmt.Sprintf("/juegos/detalle?id=%d", juegoID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// DeleteHandler: Soporta tanto DELETE /copias/{id} (HTMX) como POST /copias/delete (fallback)
func (c *CopiasController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var idStr string
	var juegoIDStr string

	if r.Method == http.MethodDelete {
		idStr = r.PathValue("id")
	} else if r.Method == http.MethodPost {
		_ = r.ParseForm()
		idStr = r.FormValue("id")
		juegoIDStr = r.FormValue("juego_id")
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	parsedID, err := strconv.ParseInt(idStr, 10, 32)
	if err == nil && parsedID > 0 {
		_ = c.service.DeleteCopia(r.Context(), int32(parsedID))
	}

	// Si es HTMX, 200 OK vacío elimina la fila directamente
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if juegoIDStr != "" {
		http.Redirect(w, r, fmt.Sprintf("/juegos/detalle?id=%s", juegoIDStr), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
