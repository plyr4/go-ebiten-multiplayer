BUILD_DATE_TAG=$(shell date "+d%y.%m.%d-t%H.%M")

IMAGE_PUBLISH_PATH='davidvader/go-ebiten-multiplayer

debug:
	@echo "running program directly"
	go run main.go

debug-wasm:
	@echo "serving wasm directly"
	./serverwasm.sh

kill-wasm:
	@echo "killing any frozen wasm processes related to port 8080"
	lsof -P | grep 8080 | grep nc | awk '{print $$2}' | xargs kill -9
	lsof -P | grep 8080 | grep wasmserve | awk '{print $$2}' | xargs kill -9
	lsof -P | grep 8080 | grep curl | awk '{print $$2}' | xargs kill -9

up:
	@echo "building image"
	docker build -t game:local .
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
