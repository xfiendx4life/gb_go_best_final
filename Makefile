test:
	go test -v ./... -coverprofile=coverage.out
check:
	golangci-lint run
