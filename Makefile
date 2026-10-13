MODULES := common infrastructure auth blog notification
PKGS    := $(addsuffix /...,$(addprefix ./,$(MODULES)))

COMPOSE := docker compose -f local/docker-compose.yml
POSTGRES_DSN ?= postgres://micro:micro@localhost:5432/micro?sslmode=disable

GOOSE_VERSION := v3.28.0
GOOSE := GOWORK=off go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) \
	-dir migrations postgres "$(POSTGRES_DSN)"

.PHONY: build vet lint test test-integration tidy clean proto \
	run-auth run-blog run-notification \
	up up-infra down down-v logs ps psql \
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
SIBLINGS := $(addprefix -droprequire=github.com/BenyaChef/micro/,$(MODULES))

tidy:
	@for m in $(MODULES); do (cd $$m && go mod tidy -e); done
	@for m in $(MODULES); do go mod edit $(SIBLINGS) $$m/go.mod; done
	go work sync

clean:
	rm -rf bin

PROTO_BASE  := infrastructure/grpcserver/grpcmodel
PROTO_AUTH  := common/auth/adapter/gateway/proto

proto:
	protoc -I $(PROTO_BASE) \
		--go_out=paths=source_relative:$(PROTO_BASE) \
		error.proto
	protoc -I $(PROTO_AUTH) -I $(PROTO_BASE) \
		--go_out=paths=source_relative:$(PROTO_AUTH) \
		--go-grpc_out=paths=source_relative:$(PROTO_AUTH) \
		auth_service.proto

run-auth:
	go run ./auth

run-blog:
	go run ./blog

run-notification:
	go run ./notification
up:
	$(COMPOSE) up -d --build --wait

up-infra:
	$(COMPOSE) up -d --wait postgres kafka

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
