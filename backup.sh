#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DIR="${BACKUP_DIR:-${SCRIPT_DIR}/backups}"
RETAIN_DAYS="${RETAIN_DAYS:-14}"
DATE=$(date +%Y-%m-%d_%H-%M-%S)

read_dotenv() {
    local key="$1"
    local default="$2"
    local value=""

    if [ -f "${SCRIPT_DIR}/.env" ]; then
        value=$(sed -n "s/^${key}=//p" "${SCRIPT_DIR}/.env" | tail -n 1)
        value="${value%\"}"
        value="${value#\"}"
        value="${value%\'}"
        value="${value#\'}"
    fi

    if [ -n "$value" ]; then
        printf '%s' "$value"
    else
        printf '%s' "$default"
    fi
}

POSTGRES_USER="${POSTGRES_USER:-$(read_dotenv POSTGRES_USER armpower)}"
POSTGRES_DB="${POSTGRES_DB:-$(read_dotenv POSTGRES_DB armpower)}"
BACKUP_FILE="${BACKUP_DIR}/${POSTGRES_DB}_${DATE}.dump"

mkdir -p "$BACKUP_DIR"

echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting custom-format backup..."
echo "Database: ${POSTGRES_DB}"
echo "User: ${POSTGRES_USER}"

docker compose -f "${SCRIPT_DIR}/docker-compose.yml" exec -T postgres \
    pg_dump \
    -U "${POSTGRES_USER}" \
    -d "${POSTGRES_DB}" \
    -F c \
    -Z 9 \
    --no-owner \
    --no-privileges \
    > "$BACKUP_FILE"

if [ ! -s "$BACKUP_FILE" ]; then
    rm -f "$BACKUP_FILE"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Backup failed: empty file"
    exit 1
fi

SIZE=$(du -sh "$BACKUP_FILE" | cut -f1)
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Backup saved: $BACKUP_FILE ($SIZE)"

# Удалить старые бэкапы старше RETAIN_DAYS дней.
DELETED=$(find "$BACKUP_DIR" \( -name "${POSTGRES_DB}_*.dump" -o -name "${POSTGRES_DB}_*.sql.gz" \) -mtime "+${RETAIN_DAYS}" -print -delete | wc -l)
if [ "$DELETED" -gt 0 ]; then
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Deleted $DELETED old backup(s) (older than ${RETAIN_DAYS} days)"
fi

echo "Restore example:"
echo "docker compose -f \"${SCRIPT_DIR}/docker-compose.yml\" exec -T postgres pg_restore -U \"${POSTGRES_USER}\" -d \"${POSTGRES_DB}\" --clean --if-exists --verbose < \"$BACKUP_FILE\""
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Done."
