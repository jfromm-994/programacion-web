#!/bin/bash
set -e

APP_PID=""

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

# Función auxiliar para comprobar código HTTP y contenido del body
assert_response() {
    local test_name="$1"
    local expected_status="$2"
    local expected_pattern="$3"
    local actual_status="$4"
    local body="$5"

    if [ "$actual_status" != "$expected_status" ]; then
        echo "❌ [FALLO] $test_name"
        echo "   Código HTTP esperado: $expected_status | Obtenido: $actual_status"
        echo "   Body recibido: $body"
        exit 1
    fi

    if [ -n "$expected_pattern" ]; then
        if ! echo "$body" | grep -q "$expected_pattern"; then
            echo "❌ [FALLO] $test_name"
            echo "   Patrón esperado en el body: '$expected_pattern'"
            echo "   Body recibido: $body"
            exit 1
        fi
    fi

    echo "✔ [OK] $test_name (Status: $actual_status)"
}

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
RAW=$(curl -s -w "\n%{http_code}" -X GET http://localhost:8080/juegos)
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "GET /juegos (Vacío)" "200" "null" "$STATUS" "$BODY"

echo ""
echo "--> Test 2: POST /juegos (Crear juego válido: Hollow Knight)"
RAW=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/juegos \
    -H "Content-Type: application/json" \
    -d '{"titulo": "Hollow Knight", "genero": "Metroidvania", "logros_totales": 63}')
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "POST /juegos (Crear)" "201" "Hollow Knight" "$STATUS" "$BODY"

JUEGO_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
if [ -z "$JUEGO_ID" ]; then
    echo "Error: No se pudo obtener el ID del juego creado."
    exit 1
fi
echo "   ID asignado: $JUEGO_ID"

echo ""
echo "--> Test 3: POST /juegos (Validación fallida: Título vacío, espera 400 Bad Request)"
RAW=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/juegos \
    -H "Content-Type: application/json" \
    -d '{"titulo": "", "genero": "Aventura", "logros_totales": 10}')
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "POST /juegos (Título vacío)" "400" "el título del juego no puede estar vacío" "$STATUS" "$BODY"

echo ""
echo "--> Test 4: GET /juegos/{id} (Obtener juego existente)"
RAW=$(curl -s -w "\n%{http_code}" -X GET "http://localhost:8080/juegos/$JUEGO_ID")
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "GET /juegos/$JUEGO_ID" "200" "Hollow Knight" "$STATUS" "$BODY"

echo ""
echo "--> Test 5: GET /juegos/99999 (Obtener juego inexistente, espera 404 Not Found)"
RAW=$(curl -s -w "\n%{http_code}" -X GET http://localhost:8080/juegos/99999)
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "GET /juegos/99999 (Inexistente)" "404" "Juego no encontrado" "$STATUS" "$BODY"

echo ""
echo "--> Test 6: PUT /juegos/{id} (Actualizar datos del juego)"
RAW=$(curl -s -w "\n%{http_code}" -X PUT "http://localhost:8080/juegos/$JUEGO_ID" \
    -H "Content-Type: application/json" \
    -d '{"titulo": "Hollow Knight: Silksong", "genero": "Metroidvania", "logros_totales": 70}')
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "PUT /juegos/$JUEGO_ID" "200" "Silksong" "$STATUS" "$BODY"

echo ""
echo "--> Test 7: POST /copias (Crear copia asociada al juego)"
RAW=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/copias \
    -H "Content-Type: application/json" \
    -d "{\"juego_id\": $JUEGO_ID, \"plataforma\": \"PC Steam\", \"horas_jugadas\": \"45.50\", \"logros_obtenidos\": 30}")
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "POST /copias (Crear)" "201" "PC Steam" "$STATUS" "$BODY"

COPIA_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
echo "   ID copia asignado: $COPIA_ID"

echo ""
echo "--> Test 8: GET /copias?juego_id={id} (Listar copias por juego)"
RAW=$(curl -s -w "\n%{http_code}" -X GET "http://localhost:8080/copias?juego_id=$JUEGO_ID")
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "GET /copias?juego_id=$JUEGO_ID" "200" "PC Steam" "$STATUS" "$BODY"

echo ""
echo "--> Test 9: DELETE /copias/{id} (Eliminar copia, espera 204 No Content)"
if [ -n "$COPIA_ID" ]; then
    RAW=$(curl -s -w "\n%{http_code}" -X DELETE "http://localhost:8080/copias/$COPIA_ID")
    STATUS=$(echo "$RAW" | tail -n1)
    BODY=$(echo "$RAW" | sed '$d')
    assert_response "DELETE /copias/$COPIA_ID" "204" "" "$STATUS" "$BODY"
fi

echo ""
echo "--> Test 10: DELETE /juegos/{id} (Eliminar juego, espera 204 No Content)"
RAW=$(curl -s -w "\n%{http_code}" -X DELETE "http://localhost:8080/juegos/$JUEGO_ID")
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "DELETE /juegos/$JUEGO_ID" "204" "" "$STATUS" "$BODY"

echo ""
echo "--> Test 11: GET /juegos/{id} (Comprobar borrado, espera 404 Not Found)"
RAW=$(curl -s -w "\n%{http_code}" -X GET "http://localhost:8080/juegos/$JUEGO_ID")
STATUS=$(echo "$RAW" | tail -n1)
BODY=$(echo "$RAW" | sed '$d')
assert_response "GET /juegos/$JUEGO_ID (Post-delete)" "404" "Juego no encontrado" "$STATUS" "$BODY"

echo ""
echo "¡TODAS LAS PRUEBAS DE LA API SE COMPLETARON EXITOSAMENTE!"