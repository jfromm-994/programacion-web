# Game Vault - TP2 (Persistencia de Datos)

Aplicación web multiplataforma para registrar y comparar progresos, logros y horas jugadas en videojuegos.

---

## Documentación de Persistencia

La capa de datos se implementa utilizando **PostgreSQL 16**, orquestada mediante **Docker Compose**, y mapeada a código Go fuertemente tipado mediante **`sqlc`**.

### Esquema de la Base de Datos (`db/schema/schema.sql`)

#### 1. Tabla `juegos`
Representa la información general e invariable de cada título.
* `id` (`SERIAL PRIMARY KEY`): Identificador único autoincremental.
* `titulo` (`VARCHAR(255) UNIQUE NOT NULL`): Nombre del juego (restricción de unicidad).
* `genero` (`VARCHAR(100) NOT NULL`): Categoría del videojuego.
* `logros_totales` (`INT NOT NULL DEFAULT 0`): Cantidad total de logros disponibles.
* `created_at` (`TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP`): Sello de tiempo de creación.

#### 2. Tabla `copias`
Representa las instancias específicas que posee el usuario en distintas plataformas.
* `id` (`SERIAL PRIMARY KEY`): Identificador único autoincremental.
* `juego_id` (`INT NOT NULL REFERENCES juegos(id) ON DELETE CASCADE`): Clave foránea con eliminación en cascada.
* `plataforma` (`VARCHAR(100) NOT NULL`): Entorno de juego (Steam, PlayStation, etc.).
* `horas_jugadas` (`NUMERIC(10, 2) NOT NULL DEFAULT 0.0`): Horas acumuladas.
* `logros_obtenidos` (`INT NOT NULL DEFAULT 0`): Logros desbloqueados.
* `created_at` (`TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP`): Sello de tiempo de creación.
* **Restricción de Unicidad Compuesta:** `CONSTRAINT unique_juego_plataforma UNIQUE (juego_id, plataforma)` que impide registrar copias duplicadas para la misma plataforma en un mismo juego.

---

## Instrucciones de Ejecución y Pruebas Automatizadas

### Requisitos previos
- [Go](https://go.dev/) (v1.21 o superior).
- [Docker Desktop](https://www.docker.com/) en ejecución.
- [sqlc](https://sqlc.dev/) instalado.

### Ejecución de Pruebas Automatizadas
Para ejecutar todo el ciclo de vida de pruebas (levantar base de datos, ejecutar consultas de prueba y limpiar contenedores/volúmenes), clona el repositorio y ejecuta:

```bash
./test.sh

El script test.sh realiza automáticamente:
1. Limpieza de contenedores y volúmenes anteriores (docker compose down -v).
2. Generación de código Go a través de sqlc generate.
3. Levantamiento del contenedor de PostgreSQL (docker compose up -d).
4. Espera activa hasta que PostgreSQL esté listo para recibir conexiones.
5. Ejecución de las pruebas unitarias con go test -v ./....
6. Tareas posteriores de limpieza y destrucción de contenedores y volúmenes persistentes.

---
