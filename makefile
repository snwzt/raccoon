build:
	@go build -o build/raccoon cmd/raccoon.go

clean:
	@rm -rf build