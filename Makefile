.PHONY: all clean test run-client space-invaders

all: space-invaders

clean:
	go clean ./...

test:
	go test -race -p 1 -timeout 2m -v ./...

VERSION=$(shell date +%Y%m%d-%H%M%S)-$(shell git rev-parse --verify --short HEAD)
GO_BUILD_FLAGS=
APP_NAME=space-invaders
BUILD_PATH=./build

run:
	go run ${GO_BUILD_FLAGS} main.go

space-invaders:
# 	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o ${BUILD_PATH}/$(APP_NAME)-linux ./
# 	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o ${BUILD_PATH}/$(APP_NAME)-windows.exe ./
	go build $(GO_BUILD_FLAGS) -o ${BUILD_PATH}/$(APP_NAME)-darwin ./
	printf ${VERSION} > ${BUILD_PATH}/version
	chmod a+x ${BUILD_PATH}/$(APP_NAME)-*

FORMAT_DIR=.
export FORMAT_DIR
format:
	# Note: prettier will fail if it does not match any files in the given directory.
	npx prettier --loglevel warn --write '$(FORMAT_DIR)/**/*.{md,yml,yaml,js,ts,json,html,css,scss,vue}' '!$(FORMAT_DIR)/vendor/**' '!$(FORMAT_DIR)/**/*_pb*' '!$(FORMAT_DIR)/**/package*.json' '!$(FORMAT_DIR)/**/dist/**' '!./partner.pilotfiber.com/**'
	goimports -w $(FORMAT_DIR)

lint:
	go vet ./...

