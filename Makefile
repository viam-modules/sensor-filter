test:
	go test -v

lint:
	golangci-lint run

module.tar.gz:
	go build -o SFexe cmd/main.go
	tar -czvf module.tar.gz SFexe meta.json
