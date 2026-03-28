# See https://googlecloudplatform.github.io/magic-modules/docs/getting-started/generate-providers/
# for a guide on provider generation.

default: build

# mm setup
ifneq ($(PRODUCT),)
  mmv1_compile=--product $(PRODUCT)
endif

ifneq ($(RESOURCE),)
  mmv1_compile += --resource $(RESOURCE)
endif

ifneq ($(OVERRIDES),)
  mmv1_compile += --overrides $(OVERRIDES)
endif

UNAME := $(shell uname)

# The inplace editing semantics are different between linux and osx.
ifeq ($(UNAME), Linux)
SED_I := -i
else
SED_I := -i '' -E
endif

SHOULD_SKIP_CLEAN := false # Default: do not skip
ifneq ($(SKIP_CLEAN),)
  ifneq ($(SKIP_CLEAN),false)
    SHOULD_SKIP_CLEAN := true
  endif
endif

terraform build provider: validate_environment clean-provider mmv1
	@echo "Provider generation process finished for $(VERSION) in $(OUTPUT_PATH)"


mmv1:
	@echo "Executing mmv1 build for $(OUTPUT_PATH)"; 
	@cd mmv1;\
		if [ "$(VERSION)" = "ga" ]; then \
			go run . --output $(OUTPUT_PATH) --version ga --no-docs $(mmv1_compile) \
			&& go run . --output $(OUTPUT_PATH) --version beta --no-code $(mmv1_compile); \
		else \
			go run . --output $(OUTPUT_PATH) --version $(VERSION) $(mmv1_compile); \
		fi

clean-provider: check_safe_build
	@if [ -n "$(PRODUCT)" ]; then \
		printf "\n\e[1;33mWARNING:\e[0m Skipping clean-provider step because PRODUCT ('$(PRODUCT)') is set.\n"; \
		printf " Ensure your downstream repository is synchronized with the Magic Modules branch\n"; \
		printf " to avoid potential build inconsistencies.\n"; \
		printf " Downstream repository (OUTPUT_PATH): %s\n\n" "$(OUTPUT_PATH)"; \
	elif [ "$(SHOULD_SKIP_CLEAN)" = "true" ]; then \
		printf "\e[1;33mINFO:\e[0m Skipping clean-provider step because SKIP_CLEAN is set to a non-false value ('$(SKIP_CLEAN)').\n"; \
	else \
		echo "Executing clean-provider in $(OUTPUT_PATH)..."; \
		( \
			cd $(OUTPUT_PATH) && \
			echo "---> Changing directory to $(OUTPUT_PATH)" && \
			if ! command -v git > /dev/null 2>&1; then \
				printf "\e[1;33mINFO:\e[0m Skipping git-based cleaning because git is not installed.\n"; \
			elif ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then \
				printf "\e[1;33mINFO:\e[0m Skipping git-based cleaning because $(OUTPUT_PATH) is not a git repository.\n"; \
			else \
				echo "---> Downloading Go module dependencies... (Ensures tools like gofmt can find relevant code)" && \
				go mod download && \
				echo "---> Finding tracked files to remove..." && \
				git ls-files | grep -v -E '(^\.git|^\.changelog|^\.travis\.yml$$|^\.golangci\.yml$$|^CHANGELOG\.md$$|^CHANGELOG_v.*\.md$$|^GNUmakefile$$|docscheck\.sh$$|^\.whitesource$$|^LICENSE$$|^CODEOWNERS$$|^README\.md$$|^\.go-version$$|^\.hashibot\.hcl$$|^go\.mod$$|^go\.sum$$|^examples|^scripts/)' | xargs -r git rm -f -q && \
				echo "---> Unstaging changes with git reset..." && \
				git reset -q && \
				echo "---> clean-provider actions finished. Changes have been unstaged."; \
			fi \
		) && echo "clean-provider target finished successfully."; \
	fi

clean-tgc:
	cd $(OUTPUT_PATH);\
		rm -rf ./tfplan2cai/testdata/templates/;\
		rm -rf ./tfplan2cai/testdata/generatedconvert/;\
		rm -rf ./tfplan2cai/converters/google/provider;\
		rm -rf ./tfplan2cai/converters/google/resources;\
		rm -rf ./cai2hcl/*;\
		find ./tfplan2cai/test/** -type f -exec git rm {} \; > /dev/null;\
		rm -rf ./pkg/*;\
		rm -rf ./test/*;\
		rm -rf ./cmd/*;\

tgc:
	cd mmv1;\
		go run . --version beta --provider tgc --output $(OUTPUT_PATH)/tfplan2cai $(mmv1_compile)\
		&& go run . --version ga --provider tgc_cai2hcl --output $(OUTPUT_PATH)/cai2hcl $(mmv1_compile)\
		&& go run . --version ga --provider tgc_next --output $(OUTPUT_PATH) $(mmv1_compile);\

tf-oics:
	cd mmv1;\
		go run . --version ga --provider oics --output $(OUTPUT_PATH) $(mmv1_compile);\

test:
	cd mmv1; \
		go test ./...

validate_environment: check_parameters check_safe_build

check_parameters:
# only print doctor script to console if there was a dependency failure detected.
	@./scripts/doctor 2>&1 > /dev/null || ./scripts/doctor
	@[ -d "$(OUTPUT_PATH)" ] || (printf "\n\e[1;31mERROR: directory '$(OUTPUT_PATH)' does not exist - ENV variable \033[0mOUTPUT_PATH\e[1;31m should be set to a provider directory. \033[0m \n\n" && exit 1);
	@[ -n "$(VERSION)" ] || (printf "\n\e[1;31mERROR: version '$(VERSION)' does not exist - ENV variable \033[0mVERSION\e[1;31m should be set to ga or beta \033[0m \n\n" && exit 1);


check_safe_build:
	@([ -f "$(OUTPUT_PATH)/go.mod" ] && head -n 1 "$(OUTPUT_PATH)/go.mod" | grep -q 'terraform') || \
		( \
			printf "\n\e[1;31mERROR: Validation failed for OUTPUT_PATH '$(OUTPUT_PATH)'.\n" && \
			printf "       Either go.mod is missing or the module name within it does not contain 'terraform'.\n" && \
			printf "       This is a safety check before cleaning/building. Halting.\033[0m\n\n" && \
			exit 1 \
		); \

doctor:
	./scripts/doctor

.PHONY: mmv1 test clean-provider validate_environment doctor
