test:
	@go test -race -coverprofile=coverage.out -covermode=atomic $(go list ./... | grep -v testutil)
	@go tool cover -func=coverage.out
