test:
	@go test -tags="integration" -race -coverprofile=profile.cov -v $(shell go list ./... | grep -vE 'cmd|mocks|testdata|testutil')
	@go tool cover -func=profile.cov | grep total
