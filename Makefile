test:
	go test -count=1 ./...

test-integration:
	go test -tags=integration -count=1 ./...

lint:
	golangci-lint run --exclude tx.Rollback
