build_dir=bin
TF_CONFIG_FILE=tf-dev-override.tfrc
TEST?=$$(go list -e ./... | grep -v github.com/GoogleCloudPlatform/terraform-google-conversion/v7/test/services)

build:
	GO111MODULE=on go build -o ./${build_dir}/tfplan2cai ./cmd/tfplan2cai
	GO111MODULE=on go build -o ./${build_dir}/tgc ./cmd/tgc

test:
	go version
	terraform --version
	./config-tf-dev-override.sh
	TF_CLI_CONFIG_FILE="$${PWD}/${TF_CONFIG_FILE}" GO111MODULE=on go test $(TEST) $(TESTARGS) -timeout 30m -short

test-local: mod-clean test

test-integration:
	go version
	terraform --version
	./config-tf-dev-override.sh
	TF_CLI_CONFIG_FILE="$${PWD}/${TF_CONFIG_FILE}" GO111MODULE=on go test -run=TestAcc $(TESTPATH) $(TESTARGS) -parallel 6 -timeout 60m -v ./...

mod-clean:
	git restore go.mod
	git restore go.sum
	go mod tidy

test-integration-local: mod-clean test-integration

test-go-licenses:
	cd .. && go version && go install github.com/google/go-licenses@latest
	$$(go env GOPATH)/bin/go-licenses check ./... --ignore github.com/dnaeon/go-vcr

run-docker:
	docker run -it \
	-v `pwd`:/terraform-google-conversion \
	-v ${GOOGLE_APPLICATION_CREDENTIALS}:/terraform-google-conversion/credentials.json \
	-w /terraform-google-conversion \
	--entrypoint=/bin/bash \
	--env TEST_PROJECT=${PROJECT_ID} \
	--env GOOGLE_APPLICATION_CREDENTIALS=/terraform-google-conversion/credentials.json \
	gcr.io/graphite-docker-images/go-plus;

release:
	./release.sh ${VERSION}

.PHONY: build test test-integration test-go-licenses run-docker release
