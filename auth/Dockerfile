FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.work go.work.sum ./
COPY common/go.mod common/go.sum ./common/
COPY infrastructure/go.mod infrastructure/go.sum ./infrastructure/
COPY auth/go.mod auth/go.sum ./auth/
COPY blog/go.mod blog/go.sum ./blog/
COPY notification/go.mod notification/go.sum ./notification/
RUN go mod download

COPY common ./common
COPY infrastructure ./infrastructure
COPY auth ./auth

RUN CGO_ENABLED=0 go build -trimpath -o /out/auth ./auth

FROM alpine:3.22

RUN adduser -D -u 10001 app
USER app

COPY --from=builder /out/auth /usr/local/bin/auth

EXPOSE 8081
EXPOSE 9091
EXPOSE 9101

ENTRYPOINT ["auth"]
