.PHONY: tools run build-image

tools:
	@go get github.com/google/wire/cmd/wire

run:
	@go run main.go

build-image:
	@read -p "Enter Image Name: " IMAGE_NAME; \
	docker build . -f ./build/Dockerfile -t "$$IMAGE_NAME"
