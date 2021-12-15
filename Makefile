test:
	go test ./... -coverprofile=coverage.out
check:
	golangci-lint run
