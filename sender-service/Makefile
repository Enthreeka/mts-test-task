proto:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/v1/message.proto
server:
	go run cmd/app/main.go


test:
	go test -v ./internal/...