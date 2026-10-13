# notification

Модуль `micro/notification`. Владелец таблицы `notifications`.

## Назначение

Читает события из Kafka и складывает уведомления в базу. Сейчас обрабатывается
одно действие — `REGISTERED` из топика `user`: на него создаётся уведомление
вида `welcome`.

REST-API у сервиса нет, только `/healthz` — он нужен docker-compose для
healthcheck. Вся работа идёт через consumer.

## Как устроена обработка

```
Kafka топик user
   └─ Consumer (группа notification-service)
      └─ UserEventListener   разбирает JSON в UserEvent
         └─ UserEventUseCase диспетчер по полю action
            └─ RegisteredHandler  пишет в notifications
```

**Идемпотентность.** Доставка at-least-once, дубли неизбежны. Вставка идёт
через `ON CONFLICT (event_id) DO NOTHING`: повторная доставка не создаёт вторую
запись и не считается ошибкой. Дедупликация по `event_id`, а не по `id` записи,
поэтому повтор с новым `id` тоже отсекается.

**Коммит оффсета** происходит только после успешной обработки. Сбой записи в
базу оффсет не двигает — сообщение придёт снова.

**Нечитаемое сообщение пропускается.** Если JSON не разбирается или не проходит
схему, сообщение логируется с уровнем warn и коммитится. Иначе одно битое
сообщение заблокировало бы партицию навсегда.

**Действие без обработчика** тоже пропускается: когда в enum появится
`DELETED`, а обработчика ещё нет, сервис не упадёт.

**Параллелизм.** Внутри партиции сообщения обрабатываются строго
последовательно — иначе нельзя безопасно двигать оффсет. Семафор
(`DefaultMaxConcurrentHandlers`, по умолчанию 8) ограничивает число
одновременных обработчиков **между** партициями: горутина на сообщение без
лимита запрещена правилом.

## Запуск

```sh
make up
make migrate
cp notification/.env.example notification/.env && set -a && . notification/.env && set +a
go run ./notification
```

## Переменные окружения

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `LOG_LEVEL` | debug / info / warn / error | `info` |
| `ENVIRONMENT` | Метка окружения в метриках | `local` |
| `REST_PORT` | Порт `/healthz` | `8083` |
| `METRICS_PORT` | Порт `/metrics` для Prometheus | `9103` |
| `POSTGRES_DSN` | Строка подключения | обязательна |
| `POSTGRES_MAX_CONNS` | Верхняя граница пула | `10` |
| `POSTGRES_MIN_CONNS` | Тёплые соединения | `2` |
| `KAFKA_BROKERS` | Брокеры через запятую | `localhost:29092` |
| `KAFKA_GROUP_ID` | Consumer group | `notification-service` |

## Контракт

Входящее сообщение — топик `user`, ключ `user_id`:

```json
{
  "event_id": "uuid",
  "action": "REGISTERED",
  "occurred_at": "2026-10-12T10:00:00Z",
  "data": { "user_id": "uuid", "email": "a@b.c" }
}
```

Запись в `notifications`: `event_id` и `user_id` из события, `kind = welcome`,
`payload` — исходное событие целиком, `sent_at` — время обработки.

| Метод | Путь | Ответ |
|---|---|---|
| GET | `/healthz` | `200 {status, service}` |

## Коды ошибок

Наружу сервис ничего не отдаёт, поэтому коды видны только в логах.

| Код | Когда |
|---|---|
| `021072fc-042` | Не создать consumer group — сервис не стартует |
| `021072fc-043` | Сессия consumer'а упала |
| `062b5dd8-001` | В слушатель пришло пустое сообщение |
| `0f96d603-001` | В юзкейс пришло пустое событие |
| `298ef509-*` | Событие не разобралось или не прошло схему (пропускается) |
| `c1195890-010` | Не удалось сохранить уведомление — оффсет не коммитится |
| `6ce2099a-010` | Сбой вставки в `notifications` |
| `6ce2099a-011` | Сбой выборки из `notifications` |
| `bd49bce2-*` | Невалидные поля уведомления при сборке сущности |

## Проверка

```sh
# регистрируем юзера в auth
curl -s -X POST localhost:8081/auth/register -H 'Content-Type: application/json' \
  -d '{"email":"a@b.c","password":"correct horse battery staple"}'

# в логах notification: "welcome notification stored"
# в таблице появилась запись
make psql
=> SELECT event_id, user_id, kind, sent_at FROM notifications ORDER BY sent_at DESC LIMIT 5;

# повторная доставка того же события дубля не создаёт:
=> SELECT event_id, count(*) FROM notifications GROUP BY event_id HAVING count(*) > 1;
   (пусто)
```
