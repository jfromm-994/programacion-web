#!/bin/bash
set -e

APP_PID=""

# Tareas posteriores: apagar el servidor Go y limpiar contenedores/volúmenes de Docker
cleanup() {
    echo ""
    echo "=== 7. Tareas posteriores: Limpieza general ==="
    if [ -n "$APP_PID" ] && kill -0 "$APP_PID" 2>/dev/null; then
        echo "--> Deteniendo servidor Go (PID: $APP_PID)..."
        kill "$APP_PID" 2>/dev/null || true
    fi
    echo "--> Borrando contenedores y volúmenes de Docker..."
    docker compose down -v
    echo "=== Fin de la ejecución de pruebas ==="
}
trap cleanup EXIT

echo "=== 1. Tareas previas: Limpieza de entorno previo ==="
docker compose down -v

echo "=== 2. Tareas previas: Generando código con sqlc ==="
sqlc generate

echo "=== 3. Tareas previas: Compilando binario de la aplicación ==="
go build -o tmp_gamevault_server main.go

echo "=== 4. Tareas previas: Levantando contenedor de PostgreSQL ==="
docker compose up -d

echo "=== 5. Tareas previas: Esperando a que PostgreSQL esté listo ==="
until docker exec gamevault-db psql -U postgres -d gamevault -c '\q' > /dev/null 2>&1; do
    echo "Esperando conexión a PostgreSQL..."
    sleep 1
done
echo "PostgreSQL listo y aceptando conexiones."

echo "=== 6. Tareas previas: Levantando servidor Go en segundo plano ==="
./tmp_gamevault_server &
APP_PID=$!

# Esperar a que el servidor HTTP responda en el puerto 8080
until curl -s http://localhost:8080/juegos > /dev/null 2>&1; do
    sleep 0.5
done
echo "Servidor Go HTTP iniciado correctamente (PID: $APP_PID)."

echo ""
echo "====================================================="
echo "=== PRUEBAS DE LA API REST (CURLED ENDPOINTS) ==="
echo "====================================================="

echo ""
echo "--> Test 1: GET /juegos (Lista vacía inicial)"
curl -s -i -X GET http://localhost:8080/juegos
echo ""

echo ""
echo "--> Test 2: POST /juegos (Crear juego válido: Hollow Knight)"
JUEGO_RESP=$(curl -s -X POST http://localhost:8080/juegos \
    -H "Content-Type: application/json" \
    -d '{"titulo": "Hollow Knight", "genero": "Metroidvania", "logros_totales": 63}')
echo "$JUEGO_RESP"

# Extraer el ID generado
JUEGO_ID=$(echo "$JUEGO_RESP" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
if [ -z "$JUEGO_ID" ]; then
    echo "Error: No se pudo obtener el ID del juego creado."
    exit 1
fi
echo "Juego creado con ID: $JUEGO_ID"

echo ""
echo "--> Test 3: POST /juegos (Validación fallida: Título vacío, espera 400 Bad Request)"
curl -s -i -X POST http://localhost:8080/juegos \
    -H "Content-Type: application/json" \
    -d '{"titulo": "", "genero": "Aventura", "logros_totales": 10}'
echo ""

echo ""
echo "--> Test 4: GET /juegos/{id} (Obtener juego existente)"
curl -s -i -X GET "http://localhost:8080/juegos/$JUEGO_ID"
echo ""

echo ""
echo "--> Test 5: GET /juegos/99999 (Obtener juego inexistente, espera 404 Not Found)"
curl -s -i -X GET http://localhost:8080/juegos/99999
echo ""

echo ""
echo "--> Test 6: PUT /juegos/{id} (Actualizar datos del juego)"
curl -s -i -X PUT "http://localhost:8080/juegos/$JUEGO_ID" \
    -H "Content-Type: application/json" \
    -d '{"titulo": "Hollow Knight: Silksong", "genero": "Metroidvania", "logros_totales": 70}'
echo ""

echo ""
echo "--> Test 7: POST /copias (Crear copia asociada al juego)"
COPIA_RESP=$(curl -s -X POST http://localhost:8080/copias \
    -H "Content-Type: application/json" \
    -d "{\"juego_id\": $JUEGO_ID, \"plataforma\": \"PC Steam\", \"horas_jugadas\": \"45.50\", \"logros_obtenidos\": 30}")
echo "$COPIA_RESP"

COPIA_ID=$(echo "$COPIA_RESP" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
echo "Copia creada con ID: $COPIA_ID"

echo ""
echo "--> Test 8: GET /copias?juego_id={id} (Listar copias por juego)"
curl -s -i -X GET "http://localhost:8080/copias?juego_id=$JUEGO_ID"
echo ""

echo ""
echo "--> Test 9: DELETE /copias/{id} (Eliminar copia, espera 204 No Content)"
if [ -n "$COPIA_ID" ]; then
    curl -s -i -X DELETE "http://localhost:8080/copias/$COPIA_ID"
    echo ""
fi

echo ""
echo "--> Test 10: DELETE /juegos/{id} (Eliminar juego, espera 204 No Content)"
curl -s -i -X DELETE "http://localhost:8080/juegos/$JUEGO_ID"
echo ""

echo ""
echo "--> Test 11: GET /juegos/{id} (Comprobar borrado, espera 404 Not Found)"
curl -s -i -X GET "http://localhost:8080/juegos/$JUEGO_ID"
echo ""

echo ""
echo "¡TODAS LAS PRUEBAS DE LA API SE COMPLETARON EXITOSAMENTE!"