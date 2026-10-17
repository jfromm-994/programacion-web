# GameVault - API REST (TP3)

Aplicación web desarrollada en Go con arquitectura por capas (Three-Tier), persistencia en PostgreSQL mediante `sqlc` y exposición de servicios RESTful.

---

## 🏛 Arquitectura del Proyecto

El backend se organiza respetando la separación de responsabilidades:
- **`logic/`**: Capa de lógica de negocio pura en Go. Contiene las validaciones de las entidades independientes de cualquier framework web o base de datos.
- **`handlers/`**: Capa de presentación / HTTP. Maneja el ciclo de vida de las solicitudes, decodificación/codificación JSON, códigos de estado HTTP y middleware de logging.
- **`db/`**: Capa de persistencia. Contiene las definiciones DDL en `db/schema/` y consultas preparadas en `db/queries/` (`juegos.sql` y `copias.sql`), generando código tipado en `db/sqlc/`.

---

## 🚀 Requisitos Previos

- [Go](https://go.dev/) (1.21 o superior)
- [Docker](https://www.docker.com/) y Docker Compose
- [sqlc](https://sqlc.dev/)
- `bash` y `curl`

---

## ⚙️ Ejecución y Pruebas Automatizadas

El proyecto incluye un script de orquestación que automatiza la preparación, levantado, prueba integral y limpieza del entorno en un solo paso:

```bash
bash test-api.sh