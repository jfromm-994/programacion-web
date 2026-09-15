#!/bin/bash
set -e

# Asegurar limpieza de contenedores y volúmenes al finalizar (incluso si falla un test)
trap "echo '=== Tareas posteriores: Borrando contenedores y volúmenes ===' && docker compose down -v" EXIT

echo "=== 1. Limpiando entorno previo ==="
docker compose down -v

echo "=== 2. Generando código Go con sqlc ==="
sqlc generate

echo "=== 3. Levantando contenedor de PostgreSQL ==="
docker compose up -d

echo "=== 4. Esperando a que PostgreSQL esté listo ==="
until docker exec gamevault-db psql -U postgres -d gamevault -c '\q' > /dev/null 2>&1; do
  sleep 1
done

echo "=== 5. Ejecutando tests del paquete testing ==="
go test -v ./...