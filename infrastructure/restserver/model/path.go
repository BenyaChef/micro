package restservermodel

// Healthcheck дёргается каждые несколько секунд, в access-логе от него
// только шум. /metrics сюда не нужен: он на отдельном сервере, мимо роутера.
const PathHealthz = "/healthz"
