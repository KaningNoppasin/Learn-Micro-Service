.PHONY: proto-gen proto-clean docker-build docker-up docker-down docker-dev vendor vendor-clean

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

# Vendor commands
vendor:
	@echo "Creating vendor directories for all services..."
	cd api-gateway && go mod vendor
	cd user-service && go mod vendor
	@echo "✓ Vendor directories created!"

vendor-clean:
	@echo "Cleaning vendor directories..."
	rm -rf api-gateway/vendor user-service/vendor
	@echo "✓ Vendor directories cleaned!"

vendor-update: vendor-clean
	@echo "Updating dependencies and recreating vendor..."
	cd api-gateway && go mod tidy && go mod vendor
	cd user-service && go mod tidy && go mod vendor
	@echo "✓ Dependencies updated and vendor recreated!"

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build
	@echo "✓ Docker images built!"

docker-up:
	@echo "Starting services..."
	docker-compose up -d
	@echo "✓ Services started!"

docker-down:
	@echo "Stopping services..."
	docker-compose down
	@echo "✓ Services stopped!"

docker-dev:
	@echo "Starting development environment..."
	docker-compose -f docker-compose.dev.yml up --build
	@echo "✓ Development environment started!"

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v
	docker system prune -f

# Combined commands
setup: proto-gen vendor
	@echo "✓ Project setup complete!"

clean: proto-clean vendor-clean docker-clean
	@echo "✓ Project cleaned!"
