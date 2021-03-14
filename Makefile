TLS_ENABLED ?= false

gen:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

clean:
	rm pb/*.go

server:
	go run cmd/server/main.go -port 5000 -tls=${TLS_ENABLED}

server1:
	go run cmd/server/main.go -port 50051 -tls=${TLS_ENABLED}

server2:
	go run cmd/server/main.go -port 50052 -tls=${TLS_ENABLED}

client:
	go run cmd/client/main.go -address 0.0.0.0:5000 -tls=${TLS_ENABLED}

test:
	go test -cover -race ./...

cert:
	cd cert; ./gen.sh; cd ..

.PHONY: gen clean server client test cert