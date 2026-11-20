# GameVault - Aplicación Web Reactiva con HTMX (TP6)

## Dominio de la Aplicación

**GameVault** es una plataforma orientada a la gestión, catalogación y seguimiento de videojuegos personales y sus copias en distintas plataformas (PC Steam, PlayStation, Nintendo Switch, Xbox, etc.).

El sistema modela una relación uno a muchos (**1:N**) entre dos entidades:
* **Juegos (`juegos`)**: Representa el título canónico del catálogo (título único, género principal y cantidad de logros totales disponibles).
* **Copias (`copias`)**: Representa las instancias o versiones físicas/digitales que posee el usuario, registrando horas jugadas y logros desbloqueados en cada plataforma.

### Métricas Agregadas en Catálogo
El catálogo principal calcula y proyecta directamente desde PostgreSQL:
* **Horas Totales**: Suma acumulada de las horas jugadas en todas las copias asociadas a un título.
* **Progreso Máximo de Logros**: Máxima cantidad de logros obtenidos entre todas las plataformas registradas para dicho juego.

---

## Arquitectura y Tecnologías (TP6)

Siguiendo las directivas de la cátedra para el TP6, la aplicación evoluciona de un modelo Server-Side Rendering tradicional con recargas completas a una experiencia de usuario fluida y reactiva tipo **Single Page Application (SPA)** mediante **HTMX**:

* **Interactividad Declarativa con HTMX:** Mediante atributos HTML como `hx-post`, `hx-delete`, `hx-target` y `hx-swap`, el cliente envía peticiones asíncronas (AJAX) y actualiza fragmentos específicos del DOM sin recargar la página.
* **Altas Reactivas y Limpieza con Doble Swap (OOB):**
  * Al enviar el formulario de alta, el servidor procesa la entidad y retorna el fragmento HTML de la tabla actualizada reemplazando `#juegos-table-container` o `#copias-table-container`.
  * En la misma respuesta HTTP, se utiliza un intercambio fuera de banda (*Out-of-Band Swap*) con `hx-swap-oob="true"` para resetear los campos del formulario de manera automática y sin JavaScript manual.
* **Bajas Inline con Confirmación Nativa:**
  * Cada fila de la tabla incluye un botón de eliminación que dispara una petición `DELETE` directa al endpoint REST (`/juegos/{id}` o `/copias/{id}`).
  * Se utiliza `hx-confirm` para interceptar la acción con un cuadro de diálogo nativo del navegador antes de realizar la petición.
  * El servidor responde con código HTTP `200 OK` y cuerpo vacío; HTMX retira automáticamente del DOM únicamente la fila correspondiente (`hx-target="closest tr"` con `hx-swap="outerHTML"`).
* **Componentes Tipados y Modulares (`views/`):** Construidos y compilados en Go mediante **`templ`**, divididos en componentes atómicos para páginas completas (`JuegosPage`, `CopiasPage`) y fragmentos reactivos (`JuegosTable`, `JuegoForm`, `CopiasTable`, `CopiaForm`).
* **Lógica de Negocio y Normalización (`logic/`):** Capa desacoplada encargada de validar reglas de dominio (evitar logros de copia superiores a los del juego base, unicidad de plataformas y campos obligatorios) y normalización canónica de cadenas (Title Case neutral con `cases.Title`).
* **Persistencia (`db/`):** PostgreSQL gestionado en contenedor Docker con consultas SQL tipadas y generadas mediante **`sqlc`**.
* **Estilos Semánticos:** Interfaz limpia y responsive utilizando **Pico.css** vía CDN.

---

## Requisitos Previos

* [Go](https://go.dev/) (1.21 o superior)
* [Docker](https://www.docker.com/) y Docker Compose
* Intérprete `bash`

*(Nota: Las herramientas de generación de código `sqlc` y `templ` se ejecutan automáticamente vía `go run` dentro de los scripts de ejecución y prueba, sin requerir instalación global previa en el sistema).*

---

## Ejecución del Proyecto

La cátedra solicita un comando único para orquestar el ciclo completo de la aplicación. Para ejecutar el proyecto, clone el repositorio, sitúese en la rama `tp6` y ejecute:

```bash
bash run-app.sh
```

El script `run-app.sh` automatiza todas las fases requeridas:
1. **Tareas previas:**
   * Limpieza de entornos y volúmenes residuales (`docker compose down -v`).
   * Generación de código SQL tipado con `sqlc generate`.
   * Compilación de plantillas HTML tipadas con `templ generate`.
   * Compilación del binario de Go.
   * Inicialización del contenedor de PostgreSQL con Docker Compose y espera activa de conexión.
2. **Ejecución de la aplicación:**
   * Inicio del servidor web en `http://localhost:8080`.
   * Precarga automática de registros iniciales (*Seed Data*) para visualización inmediata del catálogo (*Hollow Knight*, *Celeste*, *Dead Cells*).
3. **Tareas posteriores:**
   * Al interrumpir el servidor con `Ctrl + C`, se ejecuta un `trap` de limpieza que detiene el binario y elimina los contenedores y volúmenes creados.

---

## Verificación de Integridad de la Base de Datos

Para ejecutar las pruebas unitarias y de integración sobre la capa de datos:

```bash
bash test.sh
```
