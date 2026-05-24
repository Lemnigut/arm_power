# Обновление бекенда на сервере

Инструкция для проекта `arm_back`. Бекенд запускается через Docker Compose:

- `postgres` - PostgreSQL 16 с данными в Docker volume `pgdata`
- `api` - Go-сервер, собирается из `Dockerfile`
- конфигурация берется из `.env`

Для обычного обновления кода достаточно забрать изменения из Git и пересобрать только контейнер `api`. PostgreSQL останавливать не нужно.

## 1. Подключиться к серверу

```bash
ssh USER@SERVER
cd /path/to/arm_back
```

Замените `/path/to/arm_back` на реальный путь к папке проекта на сервере.

## 2. Проверить, что это нужный репозиторий

```bash
git remote -v
git branch --show-current
git status --short
docker compose ps
```

Ожидаемый репозиторий:

```text
https://github.com/Lemnigut/arm_power.git
```

Если `git status --short` показывает локальные изменения, перед обновлением нужно понять, что это за изменения. Не делайте `git reset --hard`, если не уверены, что эти изменения можно потерять.

## 3. Обновить код

git fetch origin
git checkout main
git pull --ff-only origin main
```

Если `git pull --ff-only` завершился ошибкой, значит на сервере есть локальные коммиты или ветка разошлась с `origin/main`. В этом случае лучше остановиться и отдельно разобраться с Git, чтобы не потерять изменения.

## 4. Пересобрать и перезапустить только API

```bash
docker compose up -d --build api
```

PostgreSQL при этом не пересоздается и данные не трогаются.

## 5. Проверить, что бекенд живой

```bash
docker compose ps
docker compose logs --tail=100 api
curl -fsS "http://127.0.0.1:${PORT:-8080}/health"
```

Ожидаемый ответ от `/health`:

```json
{"status":"ok"}
```

Если перед сервером стоит Nginx или другой reverse proxy, дополнительно проверьте внешний URL бекенда.

## Короткая команда

```bash
cd /path/to/arm_back
git pull --ff-only origin main
docker compose up -d --build api
docker compose logs --tail=100 api
curl -fsS "http://127.0.0.1:${PORT:-8080}/health"
```

## Откат

1. Найти предыдущий рабочий коммит:

```bash
git log --oneline -10
```

2. Переключиться на него и пересобрать API:

```bash
git checkout COMMIT_HASH
docker compose up -d --build api
curl -fsS "http://127.0.0.1:${PORT:-8080}/health"
```

## Когда нужны миграции

Обычное обновление их не требует. Миграции нужны только если в новом коде появились новые файлы в `migrations/` или приложение начало ожидать новые таблицы/колонки.

Проверить, были ли изменения в миграциях после обновления:

```bash
git diff --name-only HEAD@{1} HEAD -- migrations
```

Если команда ничего не вывела, миграции трогать не нужно.

Если миграции появились, сначала сделайте бэкап:

```bash
chmod +x ./backup.sh
./backup.sh
```

Затем запустите миграции отдельным Docker-контейнером:

```bash
set -a
. ./.env
set +a

docker run --rm \
  --network "$(basename "$PWD")_default" \
  -v "$PWD/migrations:/migrations:ro" \
  migrate/migrate \
  -path=/migrations \
  -database "$DB_URL" \
  up
```

Если Docker Compose создал сеть с другим именем, посмотрите ее так:

```bash
docker network ls
```

И подставьте правильное имя сети в `--network`.

Миграции вниз откатывать только если это действительно нужно. Откат миграций может удалить данные или колонки. Перед этим обязательно проверьте файл `migrations/*.down.sql` и убедитесь, что есть свежий бэкап.

Пример отката одной миграции:

```bash
set -a
. ./.env
set +a

docker run --rm \
  --network "$(basename "$PWD")_default" \
  -v "$PWD/migrations:/migrations:ro" \
  migrate/migrate \
  -path=/migrations \
  -database "$DB_URL" \
  down 1
```

## Первый деплой на новый сервер

```bash
git clone https://github.com/Lemnigut/arm_power.git /path/to/arm_back
cd /path/to/arm_back
cp .env.example .env
nano .env
docker compose up -d postgres
docker compose up -d --build api
```
