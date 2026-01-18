dev:
	APP_ENV=development air

build:
	go build -o ./tmp/main.exe ./src/cmd/main.go

start: build
 	APP_ENV=production ./tmp/main.exe

lint:
	golangci-lint-v2 run ./...

format:
	go fmt ./...

setup-hooks:
	git config core.hooksPath ./hooks