CHDIR?=$$(ls -d examples/*/* | grep -v ".tf")
TF_LOG?=WARN
TF_CLI_CONFIG_FILE?=.terraformrc

default: install

build:
	@go mod download
	@go build -v .

install: build
	@go install -v .
	@mv terraform-provider-deploygate /tmp

test:
	@go test -v $(TESTARGS) -cover -timeout=120s -parallel=4 ./...

testacc: install
	@TF_ACC=1 go test -v $(TESTARGS) -cover -timeout 120m ./...

docs:
	@go generate ./...

terraform: install
	@for dir in $(CHDIR); do \
		TF_LOG=$(TF_LOG) TF_CLI_CONFIG_FILE=$(TF_CLI_CONFIG_FILE) terraform -chdir=$$dir $(COMMAND); \
	done

plan:
	@COMMAND=plan make terraform

apply:
	@COMMAND=apply make terraform

destroy:
	@COMMAND=destroy make terraform

.PHONY: docs testacc test install build plan apply destroy terraform
