BUILD_DATE_TAG=$(shell date "+d%y.%m.%d-t%H.%M")

IMAGE_PUBLISH_PATH=docker.io/davidvader/go-ebiten-multiplayer

client:
	@echo "running client directly"
	go run client/main.go

server:
	@echo "running server entrypoint directly"
	./entrypoint.sh

kill-server:
	@echo "Killing any frozen wasm processes related to specified ports"
	for port in 8080 8090 8091; do \
		for process in nc wasmserve curl main; do \
			lsof -P | grep $$port | grep $$process | awk '{print $$2}' | xargs -r kill -9; \
		done; \
	done

up: build run

restart: down up

publish: build-static tag push

run:
	@echo "running container"
	docker run -d -p '8080:8080' --name=game game:local

down:
	@echo "tearing down"
	-docker kill game
	-docker rm game

build:
	@echo "building image"
	docker-buildx build -t game:local -f Dockerfile.internal .

build-static:
	@echo "building static image for linux/amd64"
	docker-buildx build --platform=linux/amd64 -f Dockerfile .

tag:
	@echo "pushing image to ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}"
	docker tag game:local ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}

push:
	@echo "pushing image to ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}"
	docker push ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}
