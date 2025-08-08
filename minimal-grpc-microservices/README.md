## Key Best Practices in This Minimal Example

- Clean Architecture: Separation between HTTP API (Gin) and gRPC service

- Proper Error Handling: Context timeouts and error responses

- Type Safety: Protocol Buffers ensure type-safe communication

- Containerization: Docker for easy deployment and scaling

- Connection Management: Single gRPC connection with proper cleanup

- Validation: Input validation at the API gateway level

- Logging: Basic logging for debugging and monitoring


## Setup
```
make setup
```

# Development
```
make docker-dev
```

# Production
```
make docker-build
make docker-up
```

# View logs
```
make docker-logs
```

# Clean up
```
make clean
```
