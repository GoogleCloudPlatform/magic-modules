

default: build

ifneq ($(PRODUCT),)
  mmv1_compile=-p products/$(PRODUCT)
  tpgtools_compile = --service $(PRODUCT)
else
  mmv1_compile=-a
  tpgtools_compile =
endif
ifneq ($(RESOURCE),)
  mmv1_compile += -t $(RESOURCE)
  tpgtools_compile += --resource $(RESOURCE)
endif
build:
	make mmv1
	make tpgtools

mmv1:
	cd mmv1;\
		bundle; \
		bundle exec compiler -e terraform -o $(OUTPUT_PATH) -v $(VERSION) $(mmv1_compile);

tpgtools:
	cd tpgtools;\
		go run . --path "api" --overrides "overrides" --output $(OUTPUT_PATH) --version $(VERSION) $(tpgtools_compile)

.PHONY: mmv1 tpgtools
