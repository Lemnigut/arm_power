#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DIR="${SCRIPT_DIR}/backups"
RETAIN_DAYS="${RETAIN_DAYS:-14}"
DATE=$(date +%Y-%m-%d)
BACKUP_FILE="${BACKUP_DIR}/armpower_${DATE}.sql.gz"

mkdir -p "$BACKUP_DIR"

echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting backup..."

docker compose -f "${SCRIPT_DIR}/docker-compose.yml" exec -T postgres \
    pg_dump -U "${POSTGRES_USER:-armpower}" "${POSTGRES_DB:-armpower}" \
    | gzip > "$BACKUP_FILE"

SIZE=$(du -sh "$BACKUP_FILE" | cut -f1)
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Backup saved: $BACKUP_FILE ($SIZE)"

# Удалить бэкапы старше RETAIN_DAYS дней
DELETED=$(find "$BACKUP_DIR" -name "armpower_*.sql.gz" -mtime "+${RETAIN_DAYS}" -print -delete | wc -l)
if [ "$DELETED" -gt 0 ]; then
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Deleted $DELETED old backup(s) (older than ${RETAIN_DAYS} days)"
fi

echo "[$(date '+%Y-%m-%d %H:%M:%S')] Done."
