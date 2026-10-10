# auth

Модуль `micro/auth`. Владелец таблицы `users`.

## Назначение

Регистрация, вход и выдача access-токена. Единственный сервис, который знает
секрет подписи JWT: остальные не разбирают токен сами, а проверяют его через
gRPC-метод `ValidateToken`.

## Запуск

```sh
make up
make migrate
cp auth/.env.example auth/.env && set -a && . auth/.env && set +a
go run ./auth
```

## Переменные окружения

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `LOG_LEVEL` | debug / info / warn / error | `info` |
| `REST_PORT` | Порт REST-API | `8081` |
| `GRPC_PORT` | Порт gRPC для межсервисных вызовов | `9091` |
| `POSTGRES_DSN` | Строка подключения | обязательна |
| `POSTGRES_MAX_CONNS` | Верхняя граница пула | `10` |
| `POSTGRES_MIN_CONNS` | Тёплые соединения | `2` |
| `JWT_SECRET` | Ключ подписи, не короче 32 символов | обязательна |
| `JWT_TTL` | Время жизни access-токена | `1h` |
| `PASSWORD_HASH_COST` | Cost bcrypt | `10` |
| `KAFKA_BROKERS` | Брокеры через запятую | `localhost:29092` |

## Контракт REST

| Метод | Путь | Ответ |
|---|---|---|
| POST | `/auth/register` | `201 {user_id, email}` |
| POST | `/auth/login` | `200 {access_token, expires_in}` |
| GET | `/auth/me` | `200 {user_id, email}`, требует `Authorization: Bearer` |
| GET | `/healthz` | `200 {status, service}` |

## Контракт gRPC

Порт 9091, сервис `micro.auth.schema.AuthService`. Контракт лежит одной копией
в `common/auth/adapter/gateway/proto/auth_service.proto`, оттуда же его берут
потребители. Готовый клиент — `common/auth/adapter/gateway/`.

```proto
rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
```

Ответ — `oneof`: либо `token.user_id`, либо `error{code, msg}`. Ошибка едет
**внутри ответа**, а не gRPC-статусом, поэтому `err` от вызова означает только
транспортный сбой. Коды в `error.code` те же, что в таблице ниже.

## События

При успешной регистрации уходит сообщение в топик `user`, ключ — `user_id`.

```json
{
  "event_id": "uuid",
  "action": "REGISTERED",
  "occurred_at": "2026-10-07T10:00:00Z",
  "data": { "user_id": "uuid", "email": "a@b.c" }
}
```

Один топик на домен, действие в поле `action`. Публикация идёт после вставки
в базу: если Kafka недоступна, ошибка попадает в лог, но регистрация
завершается успешно — клиент получит 201. При недоступной Kafka **на старте**
сервис не поднимется: продюсер обязательная зависимость.

## Коды ошибок

| Код | HTTP | Когда |
|---|---|---|
| `64c4a693-001` | 400 | Тело запроса не разобрать |
| `64c4a693-002` | 400 | Нет обязательного поля |
| `545af577-002` | 400 | Email не проходит формат |
| `d3c2ee1e-003` | 400 | Пароль короче 8 символов |
| `d3c2ee1e-031` | 401 | Неверная пара email/пароль |
| `10b754c6-002` | 401 | Токен невалиден |
| `10b754c6-003` | 401 | Токен просрочен |
| `8b8c0e30-001` | 401 | Нет заголовка `Authorization: Bearer` |
| `d3c2ee1e-020` | 404 | Пользователь не найден |
| `d3c2ee1e-030` | 409 | Email уже занят |

## Проверка

```sh
curl localhost:8081/healthz
curl -X POST localhost:8081/auth/register -H 'Content-Type: application/json' \
  -d '{"email":"a@b.c","password":"correct horse battery staple"}'
```
