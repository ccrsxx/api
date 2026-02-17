PKG ?= ./...

dev:
	APP_ENV=development air

build:
	go build -o ./tmp/main.exe ./src/cmd/api/main.go

start: build
	APP_ENV=production ./tmp/main.exe

tidy:
	go mod tidy

test:
	go test -cover $(PKG)

test-race:
	go test -cover -race $(PKG)

test-watch:
	-gotestsum -- -cover $(PKG)
	gotestsum --watch -- -cover $(PKG)

test-coverage:
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out

lint:
	golangci-lint run $(PKG)

format:
	go fmt $(PKG)

setup-hooks:
	npm i -g @commitlint/cli @commitlint/config-conventional
	git config core.hooksPath ./.hooks

setup-scripts:
	cd scripts && npm i

generate-openapi: setup-scripts
	npm run start:no-commit --prefix ./scripts

generate-openapi-commit: setup-scripts
	npm run start --prefix ./scripts