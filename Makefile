.PHONY: generate tidy test e2e_test

generate:
	go generate ./collector

tidy:
	go mod tidy

test:
	go test ./collector

build/run_server:
	go build -o build/. cmd/run_server.go

integration_test:
	go test --tags=integration ./collector

e2e_test:
	go test ./e2e_test
