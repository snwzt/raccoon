build:
	@go build -o build/raccoon cmd/raccoon.go

test:
	@go test -v ./tests/...

clean:
	@rm -rf build