# micro

Учебная монорепа: три Go-микросервиса, общающиеся через REST, gRPC и Kafka.
Архитектура и раскладка папок повторяют рабочую монорепу (Clean Architecture + DDD).

План реализации лежит локально в `docs/PLAN.md` (не в гите).

## Модули

| Каталог | Модуль | Что это |
|---|---|---|
| `common/` | `micro/common` | Доменные примитивы и межсервисные контракты |
| `infrastructure/` | `micro/infrastructure` | Логгер, конфиг, HTTP/gRPC/Kafka-обвязка |
| `auth/` | `micro/auth` | Регистрация, логин, выдача и проверка токенов |
| `blog/` | `micro/blog` | Статьи |
| `notification/` | `micro/notification` | Обработка событий из Kafka |

Модули связаны через `go.work`. Сервисы не зависят друг от друга: импортировать
код чужого сервиса нельзя технически — его модуля нет в зависимостях. Общее
живёт только в `common` и `infrastructure`.

## Как это связано

```
         REST :8081                 REST :8082            REST :8083
             │                          │                  (healthz)
          ┌──┴───┐                  ┌───┴──┐             ┌─────┴──────┐
          │ auth │◀── gRPC :9091 ───│ blog │             │notification│
          └──┬───┘   ValidateToken  └───┬──┘             └─────┬──────┘
             │                          │                      │
             │  событие REGISTERED      │                      │
             └──────▶ Kafka топик user ─┼──────────────────────┘
                                        │
                              ┌─────────┴─────────┐
                              │  Postgres: users, │
                              │ posts, notifications
                              └───────────────────┘
```

`blog` не знает секрет подписи JWT и не разбирает токен сам — на каждый
защищённый запрос спрашивает `auth` по gRPC. `posts.author_id` ссылается на
пользователя чужого сервиса **без foreign key**: каждый сервис ходит только
в свои таблицы.

## Быстрый старт

```sh
make up        # поднимает всё: postgres, kafka, kafka-ui и три сервиса
```

Миграции применяются автоматически: одноразовый контейнер `migrate` отрабатывает
до старта сервисов. Команда ждёт, пока все healthcheck'и станут зелёными.

| Сервис | Порт | Детали |
|---|---|---|
| auth | 8081 REST, 9091 gRPC | `auth/SERVICE.md` |
| blog | 8082 | `blog/SERVICE.md` |
| notification | 8083 (только healthz) | `notification/SERVICE.md` |
| postgres | 5432 | |
| kafka | 9092 внутри, 29092 с хоста | |
| kafka-ui | 8080 | |

## Сквозной сценарий

Проверяет всю цепочку: регистрация → событие в Kafka → уведомление в базе →
работа со статьями с проверкой прав.

```sh
make up

PASS='correct horse battery staple'
ALICE="alice-$RANDOM@example.com"
BOB="bob-$RANDOM@example.com"

# 1. Регистрируем двоих. На каждого auth публикует событие в топик user
curl -s -X POST localhost:8081/auth/register -H 'Content-Type: application/json' \
  -d "{\"email\":\"$ALICE\",\"password\":\"$PASS\"}"
curl -s -X POST localhost:8081/auth/register -H 'Content-Type: application/json' \
  -d "{\"email\":\"$BOB\",\"password\":\"$PASS\"}"

# 2. notification прочитал события и записал уведомления
docker compose -f deploy/docker-compose.yml logs notification | grep 'welcome notification stored'
make psql
#=> SELECT kind, user_id FROM notifications;

# 3. Берём токены
login() { curl -s -X POST localhost:8081/auth/login -H 'Content-Type: application/json' \
  -d "{\"email\":\"$1\",\"password\":\"$PASS\"}" \
  | python3 -c 'import sys,json;print(json.load(sys.stdin)["access_token"])'; }
A=$(login "$ALICE"); B=$(login "$BOB")

# 4. Алиса пишет статью
POST=$(curl -s -X POST localhost:8082/posts -H "Authorization: Bearer $A" \
  -H 'Content-Type: application/json' -d '{"title":"Статья","body":"Текст"}')
PID=$(echo "$POST" | python3 -c 'import sys,json;print(json.load(sys.stdin)["post_id"])')

# 5. Права
code() { curl -s -o /dev/null -w '%{http_code}\n' "$@"; }
code localhost:8082/posts                                               # 200 — чтение публично
code -X POST localhost:8082/posts -H 'Content-Type: application/json' \
     -d '{"title":"x","body":"y"}'                                      # 401 — без токена
code -X PUT localhost:8082/posts/$PID -H "Authorization: Bearer $B" \
     -H 'Content-Type: application/json' -d '{"title":"x","body":"y"}'  # 403 — чужая статья
code -X PUT localhost:8082/posts/$PID -H "Authorization: Bearer $A" \
     -H 'Content-Type: application/json' -d '{"title":"Правка","body":"Новое"}'  # 200 — своя
code -X DELETE localhost:8082/posts/$PID -H "Authorization: Bearer $A"  # 204
```

Событие можно посмотреть глазами в kafka-ui на <http://localhost:8080> — топик
`user`, формат JSON специально ради этого.

Повторная доставка дубля не создаёт: вставка идёт через
`ON CONFLICT (event_id) DO NOTHING`, в логе будет `event already handled`.

## Команды

```sh
make up                 # поднять всё в docker, дождаться healthcheck'ов
make down               # остановить
make down-v             # остановить и снести тома
make logs               # хвост логов всех контейнеров
make ps                 # состояние контейнеров
make psql               # psql внутрь postgres

make build              # собрать все модули в bin/
make vet                # go vet по всем модулям
make lint               # golangci-lint
make test               # go test -race, тесты с базой пропускаются
make test-integration   # то же, но с живой базой
make tidy               # go mod tidy в каждом модуле
make proto              # кодоген из .proto

make migrate            # применить миграции с хоста
make migrate-status     # состояние миграций
```

`go build ./...` из корня не работает: корень не является Go-модулем.
Используйте `make build` или заходите внутрь модуля.

## Запуск без docker

Инфраструктура всё равно нужна в docker, сервисы можно гонять с хоста:

```sh
make up && make migrate
cp auth/.env.example auth/.env && set -a && . auth/.env && set +a
go run ./auth
```

## Статус

Все 10 этапов плана закрыты. Работает сквозной сценарий выше: три сервиса в
docker, общение по REST, gRPC и Kafka, идемпотентная обработка событий,
graceful shutdown с `stop_grace_period: 30s`.

Осознанно не делалось (решения записаны в `docs/PLAN.md`): transactional
outbox, DLQ с ретраями, трейсинг, метрики, k8s-манифесты.
