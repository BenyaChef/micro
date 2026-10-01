MODULES := common infrastructure auth blog notification
PKGS    := $(addsuffix /...,$(addprefix ./,$(MODULES)))

COMPOSE := docker compose -f deploy/docker-compose.yml
POSTGRES_DSN ?= postgres://micro:micro@localhost:5432/micro?sslmode=disable

GOOSE_VERSION := v3.28.0
GOOSE := GOWORK=off go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) \
	-dir migrations postgres "$(POSTGRES_DSN)"

.PHONY: build vet lint test test-integration tidy clean \
	run-auth run-blog run-notification \
	up down down-v logs ps psql \
	migrate migrate-down migrate-status migrate-create
build:
	go build -o bin/ $(PKGS)

vet:
	go vet $(PKGS)

lint:
	golangci-lint run $(PKGS)

test:
	go test $(PKGS) -race
test-integration:
	POSTGRES_TEST_DSN="$(POSTGRES_DSN)" go test $(PKGS) -race -count=1
tidy:
	@for m in $(MODULES); do (cd $$m && go mod tidy -e); done
	go work sync

clean:
	rm -rf bin

run-auth:
	go run ./auth

run-blog:
	go run ./blog

run-notification:
	go run ./notification
up:
	$(COMPOSE) up -d --wait

down:
	$(COMPOSE) down
down-v:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

psql:
	$(COMPOSE) exec postgres psql -U micro -d micro

migrate:
	$(GOOSE) up

migrate-down:
	$(GOOSE) down

migrate-status:
	$(GOOSE) status
migrate-create:
	@test -n "$(NAME)" || (echo "укажи имя: make migrate-create NAME=add_something" && exit 1)
	$(GOOSE) create $(NAME) sql
