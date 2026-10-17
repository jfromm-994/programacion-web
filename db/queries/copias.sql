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