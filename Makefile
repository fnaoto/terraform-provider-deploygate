NAME=deploygate
BINARY=terraform-provider-${NAME}
OS_ARCH?=darwin_amd64
PROVIDER_DIR=~/.terraform.d/plugins/${OS_ARCH}

default: install

build:
	go build ${BINARY}

install: build
	mkdir -p $(PROVIDER_DIR)
	mv ${BINARY} $(PROVIDER_DIR)

test:
	go test -v $(TESTARGS) -cover -timeout=120s -parallel=4 ./...

testacc: install
	TF_ACC=1 go test -v $(TESTARGS) -cover -timeout 120m ./...

docs:
	go generate ./...

.PHONY: docs examples testacc test install build
