#!/bin/bash
GIT_TAG ?= $(shell git tag)
OS = $(shell uname)

protocinstall:
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

protoserver:
	mkdir -p pkg/authserver/v1
	protoc \
    -I./pkg/authserver \
	-I api/proto/v1 \
	-I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
	--go_out=./pkg/authserver/v1 \
	--go-grpc_out=./pkg/authserver/v1 \
	--grpc-gateway_out=logtostderr=true:./pkg/authserver/v1 \
	--openapiv2_opt json_names_for_fields=false \
	--openapiv2_opt allow_merge=true \
	--openapiv2_opt merge_file_name=v1/api \
	--openapiv2_out=logtostderr=true:./pkg/authserver \
	api/proto/v1/*.proto

ifneq ($(OS), Darwin)
	sed -i 's/version not set/$(GIT_TAG)/' pkg/authserver/v1/api.swagger.json
endif

run:
	go run cmd/authserver/main.go --config=config/config.yaml