

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
terraform build:
	make serialize
	make mmv1
	make tpgtools

mmv1:
	cd mmv1;\
		bundle; \
		bundle exec compiler -e terraform -o $(OUTPUT_PATH) -v $(VERSION) $(mmv1_compile);

tpgtools:
	cd tpgtools;\
		go run . --path "api" --overrides "overrides" --output $(OUTPUT_PATH) --version $(VERSION) $(tpgtools_compile)

validator:
	cd mmv1;\
		bundle; \
		bundle exec compiler -e terraform -f validator -o $(OUTPUT_PATH) $(mmv1_compile);

serialize:
	cd tpgtools;\
		go run . --path "api" --overrides "overrides" --mode "serialization" > temp.serial &&\
		mv -f temp.serial serialization.go;\

upgrade-dcl:
	cd tpgtools && \
		go mod edit -dropreplace=github.com/GoogleCloudPlatform/declarative-resource-client-library &&\
		go mod edit -require=github.com/GoogleCloudPlatform/declarative-resource-client-library@latest &&\
		go mod tidy;\
		MOD_LINE=$$(grep declarative-resource-client-library go.mod);\
		SUM_LINE=$$(grep declarative-resource-client-library go.sum);\
	cd ../mmv1/third_party/terraform && \
		sed -i "/declarative-resource-client-library/c$$(printf '\t')$$MOD_LINE" go.mod.erb; echo -e "$$SUM_LINE" >> go.sum



.PHONY: mmv1 tpgtools

