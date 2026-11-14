# GameVault - Aplicación Web SSR (TP5)

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

## Arquitectura y Tecnologías (TP5)

Siguiendo las directivas de la cátedra para el TP5:
* **Arquitectura:** **Server-Side Rendering (SSR)** estricto con el patrón **Post/Redirect/Get (PRG)** para el procesamiento seguro de formularios.
* **Cero JavaScript:** La aplicación no utiliza ni incluye ningún archivo `.js`. Toda la interactividad, validaciones y navegación se resuelven mediante HTML semántico nativo y estilos proporcionados por **Pico.css** vía CDN.
* **Componentes Tipados (`views/`):** Construidos y compilados en Go mediante **`templ`**.
* **Lógica de Negocio (`logic/`):** Capa desacoplada encargada de validar reglas de dominio (evitar logros de copia superiores a los del juego base, unicidad de plataformas y campos obligatorios).
* **Persistencia (`db/`):** PostgreSQL gestionado en contenedor Docker con consultas SQL tipadas y generadas mediante **`sqlc`**.

---

## Requisitos Previos

* [Go](https://go.dev/) (1.21 o superior)
* [Docker](https://www.docker.com/) y Docker Compose
* [sqlc](https://sqlc.dev/)
* [templ](https://templ.guide)
* Intérprete `bash`

---

## Ejecución del Proyecto

La cátedra solicita un comando único para orquestar el ciclo completo de la aplicación. Para ejecutar el proyecto, clone el repositorio, sitúese en la rama `tp5` y ejecute:

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
3. **Tareas posteriores:**
   * Al interrumpir el servidor con `Ctrl + C`, se ejecuta un `trap` de limpieza que detiene el binario y elimina los contenedores y volúmenes creados.

---

## Verificación de Integridad de la Base de Datos

Para ejecutar las pruebas unitarias y de integración sobre la capa de datos implementadas en entregas previas:

```bash
bash test.sh
```