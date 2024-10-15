# See https://googlecloudplatform.github.io/magic-modules/docs/getting-started/generate-providers/
# for a guide on provider generation.

default: build

# mm setup
ifeq ($(ENGINE),tpgtools)
  # we specify the product to one that doesn't
  # exist so exclusively build base tpgtools implementation
  mmv1_compile=-p does-not-exist
else ifneq ($(PRODUCT),)
  mmv1_compile=--product $(PRODUCT)
endif

# tpgtools setup
ifeq ($(ENGINE),mmv1)
  # we specify the product to one that doesn't
  # exist so exclusively build base mmv1 implementation
	tpgtools_compile = --service does-not-exist
else ifneq ($(PRODUCT),)
  tpgtools_compile = --service $(PRODUCT)
else
  tpgtools_compile =
endif

ifneq ($(RESOURCE),)
  mmv1_compile += --resource $(RESOURCE)
  tpgtools_compile += --resource $(RESOURCE)
endif

ifneq ($(OVERRIDES),)
  mmv1_compile += --overrides $(OVERRIDES)
  tpgtools_compile += --overrides $(OVERRIDES)/tpgtools/overrides --path $(OVERRIDES)/tpgtools/api
  serialize_compile = --overrides $(OVERRIDES)/tpgtools/overrides --path $(OVERRIDES)/tpgtools/api
else
  tpgtools_compile += --path "api" --overrides "overrides"
  serialize_compile = --path "api" --overrides "overrides"
endif

ifneq ($(VERBOSE),)
  tpgtools_compile += --logtostderr=1 --stderrthreshold=2
endif

UNAME := $(shell uname)

# The inplace editing semantics are different between linux and osx.
ifeq ($(UNAME), Linux)
SED_I := -i
else
SED_I := -i '' -E
endif

ifeq ($(FORCE_DCL),)
  FORCE_DCL=latest
endif
terraform build provider:
	@make validate_environment;
	make mmv1
	make tpgtools

mmv1:
	cd mmv1;\
		if [ "$(VERSION)" = "ga" ]; then \
			go run . --output $(OUTPUT_PATH) --version ga --no-docs $(mmv1_compile); \
			go run . --output $(OUTPUT_PATH) --version beta --no-code $(mmv1_compile); \
		else \
			go run . --output $(OUTPUT_PATH) --version $(VERSION) $(mmv1_compile); \
		fi

tpgtools:
	make serialize
	cd tpgtools;\
		go run . --output $(OUTPUT_PATH) --version $(VERSION) $(tpgtools_compile)

clean-provider:
	cd $(OUTPUT_PATH);\
		go mod download;\
		find . -type f -not -wholename "./.git*" -not -wholename "./.changelog*" -not -name ".travis.yml" -not -name ".golangci.yml" -not -name "CHANGELOG.md" -not -name "CHANGELOG_v*.md" -not -name "GNUmakefile" -not -name "docscheck.sh" -not -name "LICENSE" -not -name "README.md" -not -wholename "./examples*" -not -name ".go-version" -not -name ".hashibot.hcl" -print0 | xargs -0 git rm > /dev/null

clean-tgc:
	cd $(OUTPUT_PATH);\
		rm -rf ./tfplan2cai/testdata/templates/;\
		rm -rf ./tfplan2cai/testdata/generatedconvert/;\
		rm -rf ./tfplan2cai/converters/google/provider;\
		rm -rf ./tfplan2cai/converters/google/resources;\
		rm -rf ./cai2hcl/*;\
		find ./tfplan2cai/test/** -type f -exec git rm {} \; > /dev/null;\

tgc:
	cd mmv1;\
		go run . --version beta --provider tgc --output $(OUTPUT_PATH)/tfplan2cai $(mmv1_compile);\
		go run . --version beta --provider tgc_cai2hcl --output $(OUTPUT_PATH)/cai2hcl $(mmv1_compile);\

tf-oics:
	cd mmv1;\
		go run . --version ga --provider oics --output $(OUTPUT_PATH) $(mmv1_compile);\

test:
	cd mmv1; \
		go test ./...

serialize:
	cd tpgtools;\
		cp -f serialization.go.base serialization.go &&\
		go run . $(serialize_compile) --mode "serialization" > temp.serial &&\
		mv -f temp.serial serialization.go

upgrade-dcl:
	make serialize
	cd tpgtools && \
		go mod edit -dropreplace=github.com/GoogleCloudPlatform/declarative-resource-client-library &&\
		go mod edit -require=github.com/GoogleCloudPlatform/declarative-resource-client-library@$(FORCE_DCL) &&\
		go mod tidy;\
		MOD_LINE=$$(grep declarative-resource-client-library go.mod);\
		SUM_LINE=$$(grep declarative-resource-client-library go.sum);\
	cd ../mmv1/third_party/terraform && \
		sed ${SED_I} "s!.*declarative-resource-client-library.*!$$MOD_LINE!" go.mod; echo "$$SUM_LINE" >> go.sum


validate_environment:
# only print doctor script to console if there was a dependency failure detected.
	@./scripts/doctor 2>&1 > /dev/null || ./scripts/doctor
	@[ -d "$(OUTPUT_PATH)" ] || (printf " \e[1;31mdirectory '$(OUTPUT_PATH)' does not exist - ENV variable \033[0mOUTPUT_PATH\e[1;31m should be set to a provider directory. \033[0m \n" && exit 1);
	@[ -n "$(VERSION)" ] || (printf " \e[1;31mversion '$(VERSION)' does not exist - ENV variable \033[0mVERSION\e[1;31m should be set to ga or beta \033[0m \n" && exit 1);

doctor:
	./scripts/doctor

.PHONY: mmv1 tpgtools test

refresh-go:
	cd mmv1;\
		bundle exec compiler.rb -e terraform -o $(OUTPUT_PATH) -v $(VERSION) $(mmv1_compile) --go-yaml; \
		go run . --yaml --template; \
		go run . --yaml --handwritten