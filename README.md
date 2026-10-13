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
| prometheus | 9090 | |
| loki | 3100 | хранилище логов |
| grafana | 3000 | логина нет, анонимный вход |

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
docker compose -f local/docker-compose.yml logs notification | grep 'welcome notification stored'
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

## Мониторинг

`make up` поднимает Prometheus и Grafana вместе с сервисами. Дашборд и источник
данных прописаны в репозитории — настраивать руками ничего не нужно.

Открыть <http://localhost:3000>: логина нет, сразу откроется дашборд
**micro — обзор сервисов**. На нём RPS и латентность по сервисам, коды ответов
(401 и 403 видны отдельно), gRPC-вызовы в auth, публикация и обработка Kafka,
пул соединений Postgres — и **логи там же**, внизу дашборда.

Логи собирает promtail прямо из docker и складывает в Loki. Смотреть их —
в **Explore** (<http://localhost:3000/explore>), на дашборде их намеренно нет:
дашборд про цифры, Explore про текст.

Explore настроен так, что открывается сразу на Loki — он источник по умолчанию.

### Открыть с готовым запросом

```sh
make explore          # печатает ссылку, открыть её в браузере
```

Ссылка открывает Explore с уже подставленным запросом, который раскрашивает
строку по смыслу:

```
10:18:48.559  WARN   blog   request handled  POST /posts -> 401  735µs  req=160b6cc8-…
   серое      жёлтый  розов.   жирным        голубой  синий        серое   зелёный
```

| Цвет | Что им помечено |
|---|---|
| красный / жёлтый / серый | уровень: `ERROR` / `WARN` / остальное |
| розовый | сервис |
| голубой | HTTP-метод |
| синий | путь |
| зелёный | любой идентификатор: `req=`, `user=`, `event=` |
| красный в скобках | код ошибки |

Сам запрос лежит в `local/grafana/explore-query.logql`. Раскраска сделана
ANSI-кодами внутри `line_format` — других способов подсветить части строки по
смыслу Grafana не даёт.

### Что писать в Explore

```logql
{service="auth"}                                  все логи сервиса
{service=~"auth|blog|notification", level="error"} только ошибки, везде
{service="blog"} |= "403"                          поиск по подстроке
{service="notification"} | json | event_id != ""   только строки с событием
```

Строки — JSON, поэтому `| json` разбирает поля и дальше по ним можно
фильтровать: `| json | status >= 500`, `| json | took > 1s`.

### Полезное

**Размер шрифта** переключается в самом компоненте логов (включены тоглы
`newLogsPanel` и `logsPanelControls`) — там же перенос строк, сортировка и
показ времени.

**Табличный вид.** Переключатель Logs/Table над результатами: с `| json` каждое
поле становится колонкой, и ничего ни с чем не сливается без всяких цветов.

**`request_id` кликабельный.** В разобранных полях строки он подсвечен как
ссылка «Все строки этого запроса» — по клику открывается запрос, собирающий
весь путь одного обращения по всем трём сервисам. Это самый быстрый способ
разобрать, что произошло.

**Уровень в логе — это серьёзность, а не факт ошибки.** Отказ клиенту (400,
401) пишется как `warn`, сбой сервиса (5xx) — как `error`, успех — `info`.
Поэтому `level="error"` показывает настоящие аварии, а не жалобы на неверный
пароль.

**`/healthz` не логируется вообще** — иначе healthcheck раз в пять секунд на
три сервиса забивает весь поток.

## Команды

```sh
make up                 # поднять всё в docker (с пересборкой), дождаться healthcheck'ов
make up-infra           # только postgres и kafka, для локального запуска сервисов
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

make explore            # ссылка на Explore с готовым цветным запросом логов
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

## Раскладка каталогов

| Каталог | Что это |
|---|---|
| `local/` | локальное окружение: docker-compose, prometheus, grafana, loki |
| `{service}/deploy/{develop,production}/` | манифесты Kubernetes для выкатки |

Два разных смысла, которые легко перепутать: `local/` поднимает систему на твоей
машине, `{service}/deploy/` описывает, как она живёт в кластере. Раскладка
второго повторяет рабочую монорепу.

## CI

`.github/workflows/ci.yml` на каждый push и PR: сборка, `go vet`, golangci-lint,
тесты против живой базы, проверка что соседние модули не уехали в `require`,
сборка `kustomize build` всех шести оверлеев, проверка что в `secret.yaml` только
заглушки, и сборка трёх docker-образов.

## CD

`.github/workflows/release.yml` по push в `main`:

- `dorny/paths-filter` определяет, какие сервисы изменились. Правка `common/`
  или `infrastructure/` пересобирает все три
- сборка и push в GHCR с тегами `sha-<commit>` и `latest`
- выкатка в кластер — шаг есть, но **включается только когда есть куда**:
  переменная репозитория `DEPLOY_ENABLED=true` и секрет `KUBECONFIG` (base64).
  Без них джоба честно пишет, что образ собран, а кластера нет

Манифесты лежат в `{service}/deploy/{develop,production}/` и собираются
kustomize. Секреты (`POSTGRES_DSN`, `JWT_SECRET`) в репозитории — заглушки
`REPLACE_ME`; реальные значения кладутся в кластер отдельно:

```sh
kubectl -n develop create secret generic auth-micro \
  --from-literal=POSTGRES_DSN='postgres://...' \
  --from-literal=JWT_SECRET='...'
```

## Статус

Все 10 этапов плана закрыты, сверх плана добавлены мониторинг и CI. Работает сквозной сценарий выше: три сервиса в
docker, общение по REST, gRPC и Kafka, идемпотентная обработка событий,
graceful shutdown с `stop_grace_period: 30s`.

Осознанно не делалось (решения записаны в `docs/PLAN.md`): transactional
outbox, DLQ с ретраями, трейсинг, метрики, k8s-манифесты.
