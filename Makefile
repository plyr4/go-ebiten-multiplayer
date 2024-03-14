BUILD_DATE_TAG=$(shell date "+d%y.%m.%d-t%H.%M")

IMAGE_PUBLISH_PATH=davidvader/go-ebiten-multiplayer

debug:
	@echo "running program directly"
	go run main.go

debug-wasm:
	@echo "serving wasm directly"
	./serverwasm.sh

kill-wasm:
	@echo "Killing any frozen wasm processes related to specified ports"
	for port in 8080 8090 8091; do \
		for process in nc wasmserve curl main; do \
			lsof -P | grep $$port | grep $$process | awk '{print $$2}' | xargs -r kill -9; \
		done; \
	done

up: build run

build:
	@echo "building image"
	docker-buildx build -t game:local -f Dockerfile .

run:
	@echo "running container"
	docker run -d -p '8080:8080' --name=game game:local

down:
	@echo "tearing down"
	-docker kill game
	-docker rm game

restart: down up

publish:
	@echo "publishing image to ${IMAGE_PUBLISH_PATH}"
	docker-buildx build --platform=linux/amd64 -t ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG} .
	docker push ${IMAGE_PUBLISH_PATH}:${BUILD_DATE_TAG}
