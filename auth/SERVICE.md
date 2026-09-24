# auth

Модуль `micro/auth`.

## Назначение

TODO: заполнить на этапе реализации сервиса.

## Запуск

```sh
go run ./auth
```

## Переменные окружения

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `LOG_LEVEL` | Уровень логирования: debug / info / warn / error | `info` |
| `REST_PORT` | Порт REST-API | `8081` |

## Проверка

```sh
curl localhost:8081/healthz
```
