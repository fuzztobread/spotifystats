build:
	go build -o bin/spotistats cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./...

clean:
	rm -rf bin
