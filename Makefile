test:
	@echo "Running Tests"
	go test -v ./...

clean:
	rm ./bin -rf

getverifiers:
	@echo "Installing golangci-lint" && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	go mod tidy

lint:
	@echo "Running $@"
	golangci-lint run -c ../.golangci.yml

build:
	@echo "Running $@"
	@go build -ldflags=\
	"-X 'github.com/codescalersinternships/tfvpn/cmd.commit=$(shell git rev-parse HEAD)'\
	 -X 'github.com/codescalersinternships/tfvpn/cmd.version=$(shell git tag --sort=-version:refname | head -n 1)'"\
	 -o bin/tfvpn main.go
