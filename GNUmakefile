# +-------------------------------------------------+-----------------------------------------------------+
# |                                   magic-modules - build commands                                      |
# +-------------------------------------------------+-----------------------------------------------------+
# | make                                            | make build                                          |
# |                                                 |                                                     |
# |  runs make build                                |  builds the terraform provider with                 |
# |                                                 |  both mmv1 and tpgtools                             |
# +-------------------------------------------------+-----------------------------------------------------+
# | make upgrade-dcl                                | make validate_environment:                          |
# |                                                 |                                                     |
# |   a command to upgrade the version of the dcl   |   a prestep to make build to ensure that the        |
# |   used within the terraform provider and        |   build environment is suitable for a build. This   |
# |   tpgtools.                                     |   command is fallible... see the                    |
# |                                                 |   troubleshooting page(TODO) for more help          |
# +-------------------------------------------------+-----------------------------------------------------+
# |                                   magic modules - parameters                                          |
# +-------------------------------------------------------------------------------------------------------+
# | OUTPUT_PATH (required) - the location of the provider you are unioning generation into. Note          |
# |   that magic-modules templates files and then unions these generated files with your local provider   |
# |   thus its always recommended you ensure magic-modules and your local provider are up to date.        |
# +-------------------------------------------------------------------------------------------------------+
# | VERSION (required) - the version of the provider you are building into.                               |
# |   valid values are 'ga' and 'beta'.                                                                   |
# +-------------------------------------------------------------------------------------------------------+
# | ENGINE - modifies make build to build either tpgtools or mmv1.                                        |
# |   valid values are 'mmv1' or `tpgtools`. Note that running make build with ENGINE set to `tpgtools`   |
# |   will build any prerequisit files from mmv1 required for tpgtools.                                   |
# +-------------------------------------------------------------------------------------------------------+
# | PRODUCT - specifies the product to build. Equivalent to the folder name within                        |
# |   `/mmv1/products` or `/tpgtools/api`. Useful to scope the build to your changes.                     |
# |   Note that mmv1 files in `mmv1/third_party/terraform` on run of make build are always built when     |
# |   building through mmv1. So if you've only modified files here you can set the PRODUCT to one that    |
# |   doesn't exist to skip the generation of all, non-handwritten, code.                                 |
# +-------------------------------------------------------------------------------------------------------+

default: build

# mm setup
ifeq ($(ENGINE),tpgtools)
  # we specify the product to one that doesn't
  # exist so exclusively build base tpgtools implementation
  mmv1_compile=-p does-not-exist
else ifneq ($(PRODUCT),)
  mmv1_compile=-p products/$(PRODUCT)
else
  mmv1_compile=-a
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
  mmv1_compile += -t $(RESOURCE)
  tpgtools_compile += --resource $(RESOURCE)
endif

ifneq ($(OVERRIDES),)
  mmv1_compile += -r $(OVERRIDES)
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
terraform build:
	@make validate_environment;
	make mmv1
	make tpgtools

mmv1:
	cd mmv1;\
		bundle; \
		bundle exec compiler -e terraform -o $(OUTPUT_PATH) -v $(VERSION) $(mmv1_compile);

tpgtools:
	make serialize
	cd tpgtools;\
		go run . --output $(OUTPUT_PATH) --version $(VERSION) $(tpgtools_compile)

validator:
	cd mmv1;\
		bundle; \
		bundle exec compiler -e terraform -f validator -o $(OUTPUT_PATH) $(mmv1_compile);

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
		sed ${SED_I} "s!.*declarative-resource-client-library.*!$$MOD_LINE!" go.mod.erb; echo "$$SUM_LINE" >> go.sum


validate_environment:
# only print doctor script to console if there was a depenecy failure detected.
	@./scripts/doctor 2>&1 > /dev/null || ./scripts/doctor
	@[ -d "${OUTPUT_PATH}" ] || (printf " \e[1;31mdirectory '${OUTPUT_PATH}' does not exist - ENV variable \033[0mOUTPUT_PATH\e[1;31m should be set to a provider directory... \033[0m \n" && exit 1);
	@[ -d "${VERSION}" ] || (printf " \e[1;31mversion '${VERSION}' does not exist - ENV variable \033[0mVERSION\e[1;31m should be set to ga or beta \033[0m \n" && exit 1);


.PHONY: mmv1 tpgtools

