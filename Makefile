gen:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

clean:
	rm pb/*.go

server:
	go run cmd/server/main.go -port 5000

client:
	go run cmd/client/main.go -address localhost:5000

test:
	go test -cover -race ./...
