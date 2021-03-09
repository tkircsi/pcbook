# protoc --proto_path=proto --go_out=plugins=grpc:pb proto/*.proto

# protoc --go_out=pb --go_opt=paths=source_relative \
#     --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
#     proto/*.proto

# protoc proto/*.proto --go_out=plugins=grpc:.

#  protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/cpu.proto