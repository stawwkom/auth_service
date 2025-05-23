export
include .env

LOCAL_BIN := $(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

PROTOC_GEN_GO := $(LOCAL_BIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(LOCAL_BIN)/protoc-gen-go-grpc

PROTO_DIR := api/auth_v1
OUT_DIR := pkg/auth_v1
PROTO_FILE := $(PROTO_DIR)/auth.proto

# Установка инструментов в локальную папку
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/bojand/ghz/cmd/ghz@latest

# Генерация gRPC и Go-кода
# Генерация gRPC и Go-кода
generate:
	mkdir -p $(OUT_DIR)
	protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=vendor.protogen \
		--proto_path=vendor.protogen/validate \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		--validate_out lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
		--plugin=protoc-gen-validate=bin/protoc-gen-validate \
		--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
		--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
		$(PROTO_FILE)

generate-auth-api:
	mkdir -p pkg/auth_login
	protoc --proto_path api/auth_login \
	--go_out=pkg/auth_login --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth_login --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/auth_login/auth_login.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/access_v1 \
	--go_out=pkg/access_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/access_v1/access.proto

gen-cert:
	openssl genrsa -out ca.key 4096
	openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, BinaryCrewInc." -days 365 -out ca.cert
	openssl genrsa -out service.key 4096
	openssl req -new -key service.key -out service.csr -config certificate.conf
	openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial \
    		-out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext


local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v
# накатывает все не накатанные миграции
local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v
# откатывает по одной миграции
local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

test:
	go clean -testcache
	go test ./internal/service/auth/tests ./internal/api/auth/tests \
		-covermode=count \
		-coverpkg=github.com/stawwkom/auth_service/internal/service/auth,github.com/stawwkom/auth_service/internal/api/auth \
		-coverprofile=cover.out \
		-v

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/stawwkom/auth_service/internal/service/auth, github.com/stawwkom/auth_service/internal/api/auth -count 5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi

grpc-load-test:
	./bin/ghz \
		--proto api/auth_v1/auth.proto \
		-i api/auth_v1,vendor.protogen,vendor.protogen/google,vendor.protogen/validate \
		--call auth.UserAPI.Get \
		--data '{"id": 2}' \
		--rps 100 \
		--total 3000 \
		--cacert certs/ca.cert \
        --cert certs/service.pem \
        --key certs/service.key \
		localhost:50051

grpc-error-load-test:
	./bin/ghz \
		--proto api/auth_v1/auth.proto \
		-i api/auth_v1,vendor.protogen,vendor.protogen/google,vendor.protogen/validate \
		--call auth.UserAPI.Get \
		--data '{"id": 0}' \
		--rps 100 \
		--total 3000 \
		--cacert certs/ca.cert \
        --cert certs/service.pem \
        --key certs/service.key \
		localhost:50051


