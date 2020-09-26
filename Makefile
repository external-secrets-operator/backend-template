protoc:
ifeq (, $(shell which protoc))
	$(error "No protoc in $(PATH), consider installing it from https://github.com/protocolbuffers/protobuf#protocol-compiler-installation")
endif
ifeq (, $(shell which protoc-gen-go))
	go install google.golang.org/protobuf/cmd/protoc-gen-go
endif
ifeq (, $(shell which protoc-gen-go-grpc))
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
endif
ifeq (, $(shell [ -f 'generated/api/api.proto' ] && echo 1))
	mkdir -p generated/api
	wget https://raw.githubusercontent.com/external-secrets-operator/operator/main/api.proto -P generated/api/
endif
	protoc --proto_path=generated/api --go_out=generated/api --go-grpc_out=generated/api generated/api/api.proto

mod:
	go mod tidy
	go mod verify

build: protoc
	go build -o bin/backend ./.

test: protoc
	go test -v ./internal/...
