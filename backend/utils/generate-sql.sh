#!/bin/bash

set -euo pipefail

# Проверка наличия аргумента для схемы
if [[ $# -lt 1 ]]; then
    echo "Usage: $0 <path-to-schema>"
    exit 1
fi

SCHEMA_PATH="$1"

# Создаем папку gen, если её нет
mkdir -p gen
SCRIPT_DIR="./gen"

REL_PATH="$(realpath --relative-to="$SCRIPT_DIR" "$SCHEMA_PATH")"

echo "$REL_PATH"

# Генерируем sqlc.yaml напрямую с нужным содержимым
cat > gen/sqlc.yaml <<EOF
version: "2"
sql:
  - schema: "${REL_PATH}"
    queries: "../queries"
    engine: "postgresql"
    gen:
      go:
        package: "gen"
        out: "./"
        emit_interface: true
EOF

# Проверяем доступность sqlc
if ! command -v sqlc &> /dev/null; then
    echo "Error: sqlc not found! Please install it first."
    exit 1
fi

# Запускаем генерацию внутри папки gen
(cd gen && sqlc generate)

echo "SQLC generation completed successfully!"
