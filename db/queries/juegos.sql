-- name: CreateJuego :one
INSERT INTO juegos (titulo, genero, logros_totales)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetJuegoByID :one
SELECT * FROM juegos
WHERE id = $1;

-- name: GetJuegoByTitulo :one
SELECT * FROM juegos
WHERE titulo = $1;

-- name: ListJuegos :many
SELECT * FROM juegos
ORDER BY titulo ASC;

-- name: ListJuegosConMetricas :many
SELECT 
    j.id, 
    j.titulo, 
    j.genero, 
    j.logros_totales,
    COALESCE(SUM(c.horas_jugadas), 0.0)::NUMERIC(10, 2) AS horas_totales,
    COALESCE(MAX(c.logros_obtenidos), 0)::INT AS max_logros_obtenidos
FROM juegos j
LEFT JOIN copias c ON j.id = c.juego_id
GROUP BY j.id, j.titulo, j.genero, j.logros_totales
ORDER BY j.titulo ASC;

-- name: UpdateJuego :exec
UPDATE juegos
SET titulo = $2, genero = $3, logros_totales = $4
WHERE id = $1;

-- name: DeleteJuego :exec
DELETE FROM juegos
WHERE id = $1;