PROJECT_NAME := Pulumi Qeycloak Resource Provider

PACK             := qeycloak
PACKDIR          := sdk
PROJECT          := github.com/raushan606/pulumi-qeycloak
NODE_MODULE_NAME := @raushan606/qeycloak
NUGET_PKG_NAME   := Pulumi.Qeycloak

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         ?= $(shell pulumictl get version)
PROVIDER_PATH   := provider
VERSION_PATH    := ${PROVIDER_PATH}/pkg/version.Version

SCHEMA_FILE     := provider/cmd/pulumi-resource-${PACK}/schema.json
GOPATH          := $(shell go env GOPATH)

WORKING_DIR     := $(shell pwd)
TESTPARALLELISM := 4

ensure::
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd examples && go mod tidy

codegen::
	(cd provider && VERSION=${VERSION} go generate cmd/${PROVIDER}/main.go)
	(cd provider && go build -o $(WORKING_DIR)/bin/${CODEGEN} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" ${PROJECT}/${PROVIDER_PATH}/cmd/$(CODEGEN))
	$(WORKING_DIR)/bin/${CODEGEN} $(SCHEMA_FILE) --version ${VERSION} 

provider::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

provider_debug::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider::
	cd provider/pkg && go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM} ./...


nodejs_sdk:: VERSION := $(shell pulumictl get version --language javascript)
nodejs_sdk::
	rm -rf sdk/nodejs
	pulumi package gen-sdk --language nodejs $(SCHEMA_FILE)
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json

java_sdk:: PACKAGE_VERSION := $(shell pulumictl get version --language generic)
java_sdk::
	rm -rf sdk/java
	pulumi package gen-sdk --language java $(SCHEMA_FILE)
	cd sdk/java/ && \
		gradle --console=plain build

.PHONY: build
build:: provider build_sdks

.PHONY: build_sdks
build_sdks: codegen nodejs_sdk java_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint::
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done


install:: install_nodejs_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin


GO_TEST := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

test_all::
	cd provider/pkg && $(GO_TEST) ./...
	cd tests/sdk/nodejs && $(GO_TEST) ./...

install_nodejs_sdk::
	-yarn unlink --cwd $(WORKING_DIR)/sdk/nodejs/bin
	yarn link --cwd $(WORKING_DIR)/sdk/nodejs/bin

test::
	cd examples && go test -v -tags=all -timeout 2h
