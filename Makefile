dev:
	APP_ENV=development air

build:
	go build -o ./tmp/main.exe ./src/cmd/api/main.go

start: build
	APP_ENV=production ./tmp/main.exe

lint:
	golangci-lint run ./...

format:
	go fmt ./...

setup-hooks:
	npm i -g @commitlint/cli @commitlint/config-conventional
	git config core.hooksPath ./.hooks