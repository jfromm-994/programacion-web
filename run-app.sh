#!/bin/bash
set -e

APP_BIN="./tmp_gamevault_server"

cleanup() {
    echo ""
    echo "====================================================="
    echo "=== Tareas posteriores: Limpieza general          ==="
    echo "====================================================="
    if [ -f "$APP_BIN" ]; then
        echo "--> Eliminando binario temporal..."
        rm -f "$APP_BIN"
    fi
    echo "--> Deteniendo y borrando contenedores y volúmenes de Docker..."
    docker compose down -v
    echo "=== Aplicación finalizada y entorno limpio ==="
}
trap cleanup EXIT INT TERM

echo "====================================================="
echo "=== 1. Tareas previas: Limpieza de entorno previo ==="
echo "====================================================="
docker compose down -v

echo ""
echo "====================================================="
echo "=== 2. Tareas previas: Generando código           ==="
echo "====================================================="
echo "--> Ejecutando sqlc generate..."
sqlc generate

echo "--> Ejecutando templ generate..."
templ generate

echo ""
echo "====================================================="
echo "=== 3. Tareas previas: Compilando la aplicación   ==="
echo "====================================================="
go build -o "$APP_BIN" main.go
echo "Binario compilado exitosamente."

echo ""
echo "====================================================="
echo "=== 4. Tareas previas: Levantando base de datos   ==="
echo "====================================================="
docker compose up -d

echo "--> Esperando a que PostgreSQL esté listo para aceptar conexiones..."
until docker exec gamevault-db psql -U postgres -d gamevault -c '\q' > /dev/null 2>&1; do
    echo "   Esperando conexión a PostgreSQL..."
    sleep 1
done
echo "PostgreSQL está listo."

echo ""
echo "====================================================="
echo "=== 5. Ejecutando GameVault (SSR con Templ)       ==="
echo "====================================================="
echo "Servidor disponible en: http://localhost:8080"
echo "Presione [Ctrl + C] para detener el servidor y limpiar el entorno."
echo "-----------------------------------------------------"

"$APP_BIN"