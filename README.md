# Game Vault - TP1

Aplicación web multiplataforma para registrar y comparar progresos, logros y horas jugadas en videojuegos.

## Descripción del Dominio y Entidades

**Game Vault** permite gestionar el progreso consolidado de videojuegos almacenados en distintas plataformas.

### Entidades Principales

#### 1. Entidad Juego (1)
Representa la información general e invariable de un título de videojuego.
* **ID:** Identificador único.
* **Título:** Nombre del juego.
* **Género:** Categoría (Acción, RPG, Estrategia, etc.).
* **Logros Totales:** Cantidad total de logros que posee el juego.

#### 2. Entidad Copia (N)
Representa la instancia específica de un juego que posee el usuario en una plataforma determinada.
* **ID:** Identificador único de la copia.
* **Juego ID:** Referencia al juego correspondiente.
* **Plataforma:** Entorno donde se posee (Steam, PlayStation, Xbox, Epic Games, etc.).
* **Horas Jugadas:** Tiempo total acumulado en esa plataforma.
* **Logros Obtenidos:** Cantidad de logros desbloqueados en esa plataforma.

---

## Requisitos previos
- [Go](https://go.dev/) (versión 1.21 o superior).

## Instrucciones de ejecución

1. Clonar el repositorio y posicionarse en la rama `tp1`:
   ```bash
   git clone https://github.com/jfromm-994/programacion-web.git
   cd programacion-web
   git checkout tp1
2. Ejecutar el servidor HTTP con Go: go run main.go
3. Abrir el navegador e ingresar a http://localhost:8080 para visualizar la aplicación.

---
