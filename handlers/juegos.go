package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"programacion-web/logic"
	"programacion-web/views"
)

type JuegosController struct {
	service *logic.Service
}

func NewJuegosController(s *logic.Service) *JuegosController {
	return &JuegosController{service: s}
}

// CatalogoHandler: GET / (o /juegos) lista juegos, POST /juegos crea uno nuevo
func (c *JuegosController) CatalogoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/juegos" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		errorMsg := r.URL.Query().Get("error")
		c.renderCatalogo(w, r, errorMsg)
	case http.MethodPost:
		c.createJuegoPRG(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// DetalleHandler: GET /juegos/detalle?id=X
func (c *JuegosController) DetalleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	parsedID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || parsedID <= 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	juegoID := int32(parsedID)
	juego, err := c.service.GetJuegoByID(r.Context(), juegoID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	copias, err := c.service.ListCopiasByJuego(r.Context(), juegoID)
	if err != nil {
		http.Error(w, "Error al obtener copias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	errorMsg := r.URL.Query().Get("error")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = views.CopiasPage(juego, copias, errorMsg).Render(r.Context(), w)
}

// DeleteHandler: POST /juegos/delete
func (c *JuegosController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	parsedID, err := strconv.ParseInt(idStr, 10, 32)
	if err == nil && parsedID > 0 {
		_ = c.service.DeleteJuego(r.Context(), int32(parsedID))
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *JuegosController) renderCatalogo(w http.ResponseWriter, r *http.Request, errorMsg string) {
	juegos, err := c.service.ListJuegosConMetricas(r.Context())
	if err != nil {
		http.Error(w, "Error al recuperar juegos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = views.JuegosPage(juegos, errorMsg).Render(r.Context(), w)
}

func (c *JuegosController) createJuegoPRG(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/?error="+url.QueryEscape("Formulario inválido"), http.StatusSeeOther)
		return
	}

	titulo := strings.TrimSpace(r.FormValue("titulo"))
	genero := strings.TrimSpace(r.FormValue("genero"))
	logrosStr := strings.TrimSpace(r.FormValue("logros_totales"))

	var logros int32
	if logrosStr != "" {
		val, err := strconv.ParseInt(logrosStr, 10, 32)
		if err != nil {
			http.Redirect(w, r, "/?error="+url.QueryEscape("Los logros deben ser un número válido"), http.StatusSeeOther)
			return
		}
		logros = int32(val)
	}

	input := logic.JuegoInput{
		Titulo:        titulo,
		Genero:        genero,
		LogrosTotales: logros,
	}

	_, err := c.service.CreateJuego(r.Context(), input)
	if err != nil {
		http.Redirect(w, r, "/?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
