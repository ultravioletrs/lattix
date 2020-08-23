BUILD_DIR = build

define compile
	go build -mod=vendor -ldflags "-s -w" -o ${BUILD_DIR}/$(1) cmd/$(1)/main.go
endef

client:
	$(call compile,client)

server:
	$(call compile,server)

clean:
	rm -rf ${BUILD_DIR}

proto:
	protoc  -I ./  ./fhesrv.proto --go_out=plugins=grpc:./

docker:
	docker build --no-cache --tag=fhe/server -f docker/Dockerfile .

runserver:
	docker-compose -f ./docker/docker-compose.yml up

.PHONY: docker
