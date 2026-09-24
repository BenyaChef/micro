# blog

Модуль `micro/blog`.

## Назначение

TODO: заполнить на этапе реализации сервиса.

## Запуск

```sh
go run ./blog
```

## Переменные окружения

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `LOG_LEVEL` | Уровень логирования: debug / info / warn / error | `info` |
| `REST_PORT` | Порт REST-API | `8082` |

## Проверка

```sh
curl localhost:8082/healthz
```
