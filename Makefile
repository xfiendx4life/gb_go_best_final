test:
	go test ./... -coverprofile=coverage.out
check:
	golangci-lint run
run:
	echo "Last commit hash" is && git rev-parse --verify HEAD && go run ./cmd/csv_processor/main.go -config="config.yaml"
