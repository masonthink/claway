#!/bin/bash
# Daily PostgreSQL backup for Claway
# Usage: Run via cron on VPS
#   0 3 * * * /opt/claway/scripts/backup-db.sh >> /var/log/claway-backup.log 2>&1

set -euo pipefail

BACKUP_DIR="/opt/claway/backups"
RETENTION_DAYS=30
CONTAINER_NAME="claway-postgres"
DB_NAME="claway"
DB_USER="claway"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/claway-${TIMESTAMP}.sql.gz"

mkdir -p "${BACKUP_DIR}"

echo "[$(date)] Starting backup..."

# Dump and compress
docker exec "${CONTAINER_NAME}" pg_dump -U "${DB_USER}" "${DB_NAME}" | gzip > "${BACKUP_FILE}"

FILESIZE=$(du -h "${BACKUP_FILE}" | cut -f1)
echo "[$(date)] Backup complete: ${BACKUP_FILE} (${FILESIZE})"

# Remove backups older than retention period
DELETED=$(find "${BACKUP_DIR}" -name "claway-*.sql.gz" -mtime +${RETENTION_DAYS} -delete -print | wc -l)
if [ "${DELETED}" -gt 0 ]; then
  echo "[$(date)] Cleaned up ${DELETED} old backup(s)"
fi

echo "[$(date)] Done. Active backups: $(ls ${BACKUP_DIR}/claway-*.sql.gz 2>/dev/null | wc -l)"
