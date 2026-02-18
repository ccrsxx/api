PKG ?= ./...

AIR_VERSION = v1.64.5
NILAWAY_VERSION = v0.0.0-20260213150243-937701de96c7
GOTESTSUM_VERSION = v1.13.0
GOLANGCI_LINT_VERSION = v2.10.1

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

lint-nill:
	nilaway $(PKG)

format:
	go fmt $(PKG)

setup-tools:
	go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)
	go install github.com/air-verse/air@$(AIR_VERSION)
	go install go.uber.org/nilaway/cmd/nilaway@$(NILAWAY_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

setup-hooks:
	npm i -g @commitlint/cli @commitlint/config-conventional
	git config core.hooksPath ./.hooks

setup-scripts:
	cd scripts && npm i

generate-openapi:
	npm run start:no-commit --prefix ./scripts

generate-openapi-commit:
	npm run start --prefix ./scripts