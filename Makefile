test-service:
	go test -v ./internal/service
check-coverage:
	go test -coverpkg=./... ./...