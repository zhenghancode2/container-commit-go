IMAGE_REF := container-commit:v0.1.0
APP_NAME := container-commit
SRC := ./...
BUILD_DIR := bin

.PHONY: all build run test fmt vet lint clean

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

test:
	go test -v $(SRC)

clean:
	rm -rf $(BUILD_DIR)

docker: build
	docker build -t $(IMAGE_REF) .