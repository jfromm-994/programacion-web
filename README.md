# GameVault - API REST (TP3)

## 🎮 Dominio de la Aplicación

**GameVault** es una plataforma diseñada para la gestión, seguimiento y catalogación de videojuegos personales y sus respectivas instancias de juego (copias físicas o digitales). 

El sistema modela una relación uno a muchos entre dos entidades principales:
* **Juegos (`juegos`)**: Representa la entidad canónica o título del catálogo general. Almacena metadatos como el título del juego, su género principal y el total de logros disponibles.
* **Copias (`copias`)**: Representa una copia específica que posee el usuario de dicho título en una plataforma concreta (PC Steam, PlayStation, Nintendo Switch, Xbox, etc.). Registra métricas de progreso particulares como las horas jugadas y la cantidad de logros desbloqueados por el usuario.

El modelo asegura integridad referencial: la eliminación de un juego base elimina en cascada todas las copias asociadas en las distintas plataformas.

---

## 🏛 Arquitectura en Tres Capas (Three-Tier)

El backend sigue un desacoplamiento estricto de responsabilidades:
* **Presentación / Transporte (`handlers/`)**: Atiende las peticiones HTTP, decodifica/codifica payloads JSON, gestiona los códigos de estado HTTP (`200`, `201`, `204`, `400`, `404`) y aplica middlewares de observabilidad. No accede directamente a la base de datos.
* **Lógica de Negocio (`logic/`)**: Capa intermedia encargada de las reglas de dominio, validaciones estructurales (campos obligatorios, restricciones numéricas) y la orquestación de operaciones hacia el almacenamiento.
* **Persistencia (`db/`)**: Gestiona el acceso a datos sobre PostgreSQL. Define el esquema DDL en `db/schema/` y consultas preparadas en `db/queries/`, utilizando `sqlc` para generar código Go tipado en `db/sqlc/`.

---

## 🚀 Requisitos Previos

- [Go](https://go.dev/) (1.21 o superior)
- [Docker](https://www.docker.com/) y Docker Compose
- [sqlc](https://sqlc.dev/)
- `bash` y `curl`

---

## ⚙️ Ejecución y Pruebas Automatizadas

El proyecto incluye un script de integración que automatiza la regeneración de código con `sqlc`, compilación, levantado del contenedor de base de datos en Docker, ejecución de las pruebas con aserciones sobre códigos HTTP y contenido del body, y la posterior limpieza:

```bash
bash test-api.sh