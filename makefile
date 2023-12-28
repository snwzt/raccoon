build:
	@go build -o bin/raccoon cmd/raccoon.go

test:
	@go test -v ./tests/...

clean:
	@rm -rf bin