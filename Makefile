BUILD_DATE_TAG=$(shell date "+d%y.%m.%d-t%H.%M")

IMAGE_PUBLISH_PATH?=docker.io/davidvader/go-ebiten-multiplayer

clt:
	@echo "running client directly"
		go run cmd/client/main.go

clt-local:
	@echo "running client directly"
	CLIENT_MULTIPLAYER=false \
		go run cmd/client/main.go

srv:
	@echo "running client directly"
		go run cmd/server/main.go

kill-srv:
	@echo "Killing any frozen wasm processes related to specified ports"
	for port in 8080 8090 8091; do \
		for process in nc wasmserve curl main; do \
			lsof -P | grep $$port | grep $$process | awk '{print $$2}' | xargs -r kill -9; \
		done; \
	done

up: build run

restart: down up

publish:
	build-static tag push

run:
	@echo "running container"
	docker run -d \
	    -p '8080:8080' \
	    -e SERVER_WS_HOST=localhost:8091 \
	    -e CLIENT_WS_HOST=localhost:8080 \
	    --name=game game:local

down:
	@echo "tearing down"
	-docker kill game
	-docker rm game

build:
	@echo "building image"
	docker-buildx build -t game:local -f Dockerfile .

build-static:
	@echo "building static image for linux/amd64"
	docker-buildx build -t game:local -f Dockerfile --platform=linux/amd64 .

tag:
	@echo "pushing image to ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}"
	docker tag game:local ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}

push:
	@echo "pushing image to ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}"
	docker push ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}
