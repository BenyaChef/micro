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
COPY blog ./blog

RUN CGO_ENABLED=0 go build -trimpath -o /out/blog ./blog

FROM alpine:3.22

RUN adduser -D -u 10001 app
USER app

COPY --from=builder /out/blog /usr/local/bin/blog

EXPOSE 8082
EXPOSE 9102

ENTRYPOINT ["blog"]
