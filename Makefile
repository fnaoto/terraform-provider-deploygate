default: install

build:
	@go mod download
	@go build -v .

test:
	@go test -v $(TESTARGS) -cover -timeout=120s -parallel=4 ./...

testacc: build
	@TF_ACC=1 go test -v $(TESTARGS) -cover -timeout 120m ./...

docs:
	@go generate ./...

.PHONY: docs examples testacc test install build
