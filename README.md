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

Модули связаны через `go.work`. Сервисы не зависят друг от друга:
в `go.mod` каждого — только `micro/common` и `micro/infrastructure`.

## Команды

```sh
make build              # собрать все модули в bin/
make vet                # go vet по всем модулям
make test               # go test -race по всем модулям
make tidy               # go mod tidy в каждом модуле
make run-auth           # запустить сервис локально
```

`go build ./...` из корня не работает: корень не является Go-модулем.
Используйте `make build` или заходите внутрь модуля.

## Локальный запуск

```sh
make build
./bin/auth &
curl localhost:8081/healthz     # {"status":"ok","service":"auth"}
```

| Сервис | Порт | Переменные |
|---|---|---|
| auth | 8081 | см. `auth/.env.example` |
| blog | 8082 | см. `blog/.env.example` |
| notification | 8083 | см. `notification/.env.example` |

## Статус

Этап 2 из 10. Готов каркас монорепы и общая инфраструктура: логгер, реестр
переменных окружения, HTTP-сервер на chi с graceful shutdown и `/healthz`.
Бизнес-логики пока нет.
