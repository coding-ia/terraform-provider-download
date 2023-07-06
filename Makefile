GOFMT_FILES?=$$(find . -name '*.go')
export GO111MODULE=on
export TF_ACC_TERRAFORM_VERSION=0.15.4
export TESTARGS=-race -coverprofile=coverage.txt -covermode=atomic

default: build

build:
	go install

dist:
	goreleaser build --single-target --skip-validate --clean

testacc:
	TF_ACC=1 go test ./internal/provider -v $(TESTARGS) -timeout 120m -count=1

httpd-start:
	@bash scripts/start-httpd.sh

httpd-stop:
	@bash scripts/stop-httpd.sh

fmt:
	gofmt -w $(GOFMT_FILES)
