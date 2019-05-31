GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=.:${GOPATH}/src --micro_out=. --go_out=. proto/image/image.proto

.PHONY: build
build: proto

	go build -o image-srv main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t microhq/image-srv:latest
