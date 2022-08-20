test:
	go test -count=1 ./...

lint:
	golangci-lint run --exclude tx.Rollback
