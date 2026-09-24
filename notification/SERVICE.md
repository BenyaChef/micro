# notification

Модуль `micro/notification`.

## Назначение

TODO: заполнить на этапе реализации сервиса.

## Запуск

```sh
go run ./notification
```

## Переменные окружения

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `LOG_LEVEL` | Уровень логирования: debug / info / warn / error | `info` |
| `REST_PORT` | Порт REST-API | `8083` |

## Проверка

```sh
curl localhost:8083/healthz
```
