test:
	go test -v

lint:
	golangci-lint run

module.tar.gz:
	go build -o module.tar.gz cmd/main.go
