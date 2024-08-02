build:
	@go build -o bin/E-commerce_REST-API cmd/main.go

test:
	@go test -v ./...

run:
	@./bin/E-commerce_REST-API