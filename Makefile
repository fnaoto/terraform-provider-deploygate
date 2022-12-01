TEST?=$$(go list ./... | grep -v 'vendor')
NAME=deploygate
BINARY=terraform-provider-${NAME}
VERSION=0.1
OS_ARCH?=darwin_amd64
EXAMPLES=examples

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${OS_ARCH}

test: 
	go test $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: install
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

examples: install
	terraform init $(EXAMPLES)
	terraform plan $(EXAMPLES)
	terraform apply $(EXAMPLES)
	terraform destroy $(EXAMPLES)

docs:
	go generate ./...

.PHONY: docs examples testacc test install build release
