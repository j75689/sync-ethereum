.PHONY: tools run build-image

tools:
	@go get github.com/google/wire/cmd/wire

build-mac-m1:
	@go build -tags dynamic

build-image:
	@read -p "Enter Image Name: " IMAGE_NAME; \
	docker build . -f ./build/Dockerfile -t "$$IMAGE_NAME"

docker-compose-build:
	@docker-compose -f ./deployment/docker-compose/docker-compose.yaml build

docker-compose-up:
	@docker-compose -f ./deployment/docker-compose/docker-compose.yaml up -d

docker-compose-down:
	@docker-compose -f ./deployment/docker-compose/docker-compose.yaml down
