MODULES := common infrastructure auth blog notification
PKGS    := $(addsuffix /...,$(addprefix ./,$(MODULES)))

.PHONY: build vet lint test tidy clean run-auth run-blog run-notification

# go build ./auth/... из корня падает: имя бинарника совпадает с именем
# каталога. -o bin/ складывает все бинарники в один каталог и снимает конфликт.
build:
	go build -o bin/ $(PKGS)

vet:
	go vet $(PKGS)

lint:
	golangci-lint run $(PKGS)

test:
	go test $(PKGS) -race

# go mod tidy работает в рамках одного модуля, поэтому обходим их циклом.
tidy:
	@for m in $(MODULES); do (cd $$m && go mod tidy); done

clean:
	rm -rf bin

run-auth:
	go run ./auth

run-blog:
	go run ./blog

run-notification:
	go run ./notification
