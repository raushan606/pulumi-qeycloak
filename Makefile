PROJECT_NAME := Pulumi Qeycloak Resource Provider

PACK             := qeycloak
PACKDIR          := sdk
PROJECT          := github.com/raushan606/pulumi-qeycloak
NODE_MODULE_NAME := @raushan606/qeycloak
NUGET_PKG_NAME   := Pulumi.Qeycloak

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         := 0.0.1
# VERSION         ?= $(shell pulumictl get version)
PROVIDER_PATH   := provider
VERSION_PATH    := ${PROVIDER_PATH}/pkg/version.Version

# Deprecated: schema.json no longer used; SDKs are generated from the provider binary
# SCHEMA_FILE     := provider/cmd/pulumi-resource-${PACK}/schema.json
GOPATH          := $(shell go env GOPATH)

WORKING_DIR     := $(shell pwd)
TESTPARALLELISM := 4

ensure::
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd examples && go mod tidy

# Deprecated: no separate codegen step; schema is served by the provider binary
# codegen::
# 	(cd provider && VERSION=${VERSION} go generate cmd/${PROVIDER}/main.go)
# 	(cd provider && go build -o $(WORKING_DIR)/bin/${CODEGEN} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" ${PROJECT}/${PROVIDER_PATH}/cmd/$(CODEGEN))
# 	$(WORKING_DIR)/bin/${CODEGEN} $(SCHEMA_FILE) --version ${VERSION} 

provider::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

provider_debug::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider::
	cd provider/pkg && go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM} ./...


nodejs_sdk:: VERSION := 0.0.1
#  VERSION := $(shell pulumictl get version --language javascript) --- IGNORE
nodejs_sdk::
	rm -rf sdk/nodejs
	pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) --language nodejs
	cd ${PACKDIR}/nodejs/ && \
		npm install && \
		npm run build
	cp README.md ${PACKDIR}/nodejs/bin/
	@if [ -f LICENSE ]; then cp LICENSE ${PACKDIR}/nodejs/bin/; fi
	cp ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/bin/
	@if [ -f ${PACKDIR}/nodejs/package-lock.json ]; then cp ${PACKDIR}/nodejs/package-lock.json ${PACKDIR}/nodejs/bin/; fi
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json

java_sdk:: PACKAGE_VERSION := 0.0.1
#  PACKAGE_VERSION := $(shell pulumictl get version --language generic)
java_sdk::
	rm -rf sdk/java
	pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) --language java
	# Copy Maven POM template to the generated SDK
	cp assets/java/pom.xml sdk/java/pom.xml
	# Attempt Maven build if mvn exists; otherwise fall back to Gradle if available.
	@if [ -f sdk/java/pom.xml ]; then \
		if command -v mvn >/dev/null 2>&1; then \
			cd sdk/java && mvn -q -Dproject.version=$(PACKAGE_VERSION) package; \
		else \
			echo "Maven not found, skipping Maven build"; \
		fi; \
	fi
	@if [ -f sdk/java/build.gradle ]; then \
		if command -v gradle >/dev/null 2>&1; then \
			cd sdk/java && gradle --console=plain build; \
		else \
			echo "Gradle not found, skipped Gradle build"; \
		fi; \
	fi

.PHONY: build
build:: provider build_sdks

.PHONY: build_sdks
build_sdks: nodejs_sdk java_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint::
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done


install:: install_nodejs_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin

# Install the provider plugin into Pulumi's local plugin cache for this VERSION
.PHONY: install_plugin_local
install_plugin_local:: provider
	@echo "Installing Pulumi plugin locally for version $(VERSION) ..."
	@mkdir -p $$HOME/.pulumi/plugins/resource-$(PACK)-v$(VERSION)
	cp $(WORKING_DIR)/bin/$(PROVIDER) $$HOME/.pulumi/plugins/resource-$(PACK)-v$(VERSION)/
	@echo "✅ Installed to $$HOME/.pulumi/plugins/resource-$(PACK)-v$(VERSION)/$(PROVIDER)"
	@echo "Installing $(PROVIDER_BINARY) to /usr/local/bin..."
	sudo cp $(WORKING_DIR)/bin/$(PROVIDER) /usr/local/bin/
	@echo "✅ Installed $(PROVIDER_BINARY)"

.PHONY: full-setup
full-setup: build install build_sdks install_plugin_local

GO_TEST := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

test_all::
	cd provider/pkg && $(GO_TEST) ./...
	cd tests/sdk/nodejs && $(GO_TEST) ./...

install_nodejs_sdk::
	- npm unlink -g ${NODE_MODULE_NAME} 2>/dev/null || true
	cd $(WORKING_DIR)/sdk/nodejs/bin && npm link

test::
	cd examples && go test -v -tags=all -timeout 2h

