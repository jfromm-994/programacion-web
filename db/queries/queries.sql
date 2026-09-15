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

-- name: UpdateJuego :exec
UPDATE juegos
SET titulo = $2, genero = $3, logros_totales = $4
WHERE id = $1;

-- name: DeleteJuego :exec
DELETE FROM juegos
WHERE id = $1;

-- name: CreateCopia :one
INSERT INTO copias (juego_id, plataforma, horas_jugadas, logros_obtenidos)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListCopiasByJuego :many
SELECT * FROM copias
WHERE juego_id = $1
ORDER BY plataforma ASC;

-- name: DeleteCopia :exec
DELETE FROM copias
WHERE id = $1;