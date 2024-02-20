.PHONY: gen clean

gen:
	@echo "Generating proto code..."
	protoc -I=pb \
   --go_out=pb --go_opt=paths=source_relative \
   --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
   --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
   bookstore.proto
	@echo "Done!"

clean:
	@echo "Cleaning up..."
	find ./pb -type f -name "*.pb.go" -delete
	find ./pb -type f -name "*.pb.gw.go" -delete
	@echo "Done!"