.PHONY: proto-gen proto-clean

proto-gen:
	@echo "Generating protobuf files..."
	@mkdir -p api-gateway/pb user-service/pb
	protoc \
		--proto_path=proto \
		--go_out=./api-gateway/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./api-gateway/pb \
		--go-grpc_opt=paths=source_relative \
		proto/user/*.proto
	protoc \
		--proto_path=proto \
		--go_out=./user-service/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./user-service/pb \
		--go-grpc_opt=paths=source_relative \
		proto/user/*.proto
	@echo "✓ Protobuf generation complete!"

proto-clean:
	rm -rf api-gateway/pb/* user-service/pb/*
