LOCAL_BIN := $(CURDIR)/bin
BUF_TAG := v1.42.0
BUF_BIN := $(LOCAL_BIN)/buf

bin: export GOBIN := $(LOCAL_BIN)
bin:
	$(info Installing binary deps...)

	go install github.com/bufbuild/buf/cmd/buf@$(BUF_TAG)
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

generate: bin
	$(BUF_BIN) generate

.PHONY: \
	# bin \ # uncomment to regenerate bin file (in RU, proxy is needed)
	generate 




