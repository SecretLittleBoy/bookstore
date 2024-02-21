.PHONY: gen clean run

PROTO_GENED_FILES = ./pb/*.pb.go ./pb/*.pb.gw.go
GO_FILES = main.go bookstore.go data.go

gen: $(PROTO_GENED_FILES)
	@echo "Generating proto code..."
	protoc -I=pb \
   --go_out=pb --go_opt=paths=source_relative \
   --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
   --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
   bookstore.proto
	@echo "Done!"

run: gen
	go run $(GO_FILES)

clean:
	@echo "Cleaning up..."
	find ./pb -type f -name "*.pb.go" -delete
	find ./pb -type f -name "*.pb.gw.go" -delete
	@echo "Done!"