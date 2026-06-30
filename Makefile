BIN      := dirwalker
CMD      := ./cmd/dirwalker
BUILD_DIR:= ./build

.PHONY: all build clean test lint run install

all: build

build:
	go build -o $(BUILD_DIR)/$(BIN) $(CMD)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

lint:
	golangci-lint run ./...

run: build
	$(BUILD_DIR)/$(BIN)

install:
	go install $(CMD)
