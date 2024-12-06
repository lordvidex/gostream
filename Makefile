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

# Вендоринг внешних proto файлов
vendor-proto: .vendor-rm vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/protoc-gen-openapiv2/options vendor-proto/validate

.vendor-rm:
	rm -rf vendor-proto
# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protovalidate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/buf/validate/
		mv vendor-proto/tmp/proto/protovalidate/buf/validate/ vendor-proto/buf/
		rm -rf vendor-proto/tmp

generate: bin vendor-proto fast-generate
fast-generate: .protoc
	# $(BUF_BIN) mod update # (not working in RU)
	# $(BUF_BIN) generate

OUT_DIR := ./pkg/api/gostream
PROTO_FILES := $(shell find api -name '*.proto')

.protoc:
	protoc \
		--proto_path=api \
		--proto_path=vendor-proto \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go_out=paths=source_relative:$(OUT_DIR) \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		--go-grpc_out=paths=source_relative:$(OUT_DIR) \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:$(OUT_DIR) \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
		--openapiv2_out=generate_unbound_methods=true:$(OUT_DIR) \
		$(PROTO_FILES)
		

build:
	go build -o $(LOCAL_BIN)/gostream ./cmd/gostream

.PHONY: \
	bin \ # uncomment to regenerate bin file (in RU, proxy is needed)
	generate \
	fast-generate \
	vendor-proto \
	.protoc \
	build







