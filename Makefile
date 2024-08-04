PHONY: gen-mocks
gen-mocks:
	go generate ./...

.PHONY: test
test:
	go test -v ./... -cover