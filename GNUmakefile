

default: build

# mm setup
ifeq ($(ENGINE),tpgtools)
  mmv1_compile=-p blahblah
else ifne ($(PRODUCT),)
  mmv1_compile=-p products/$(PRODUCT)
else
  mmv1_compile=-a
endif

# tpgtools setup
ifeq ($(ENGINE),mmv1)
	tpgtools_compile = --service blahblah
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

.PHONY: mmv1 tpgtools
