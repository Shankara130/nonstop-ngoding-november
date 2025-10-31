run:
	go run ./cmd/dashboard/main.go

build:
	go build -o bin/dashboard ./cmd/dashboard

test:
	go test ./...

docker-build:
	docker build -t nonstop-ngoding-november -f build/package/Dockerfile .
