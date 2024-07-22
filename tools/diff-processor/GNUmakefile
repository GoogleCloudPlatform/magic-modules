clone:
ifndef OWNER_REPO
	$(error OWNER_REPO is undefined)
endif
	rm -rf ./old ./new
ifdef DEPTH
	git clone --depth=$(DEPTH) https://github.com/$(OWNER_REPO).git old
else
	git clone https://github.com/$(OWNER_REPO).git old
endif
	cp -r old/ new/

clean:
	cd old/; \
	    rm -rf google/ google-beta; \
		git checkout -- .; \
		git clean -f .
	cd new/; \
	    rm -rf google/ google-beta; \
		git checkout -- .; \
		git clean -f .
	rm -rf bin/

build:
ifndef OLD_REF
	$(error OLD_REF is undefined)
endif
ifndef NEW_REF
	$(error NEW_REF is undefined)
endif
	$(MAKE) clean
	cd old/; \
		git checkout $(OLD_REF); \
		real_package_name=github.com/hashicorp/terraform-provider-google; \
		real_folder_name=google; \
		fake_package_name=google/provider/old; \
		if [ -d "google-beta" ]; then \
			mv google-beta google; \
			real_package_name=github.com/hashicorp/terraform-provider-google-beta; \
			real_folder_name=google-beta; \
		fi; \
		find . -type f -name "*.go" -exec sed -i.bak "s~$$real_package_name/$$real_folder_name~$$fake_package_name/google~g" {} +; \
		sed -i.bak "s|$$real_package_name|$$fake_package_name|g" go.mod
	cd new/; \
		git checkout $(NEW_REF); \
		real_package_name=github.com/hashicorp/terraform-provider-google; \
		real_folder_name=google; \
		fake_package_name=google/provider/new; \
		if [ -d "google-beta" ]; then \
			mv google-beta google; \
			real_package_name=github.com/hashicorp/terraform-provider-google-beta; \
			real_folder_name=google-beta; \
		fi; \
		find . -type f -name "*.go" -exec sed -i.bak "s~$$real_package_name/$$real_folder_name~$$fake_package_name/google~g" {} +; \
		sed -i.bak "s|$$real_package_name|$$fake_package_name|g" go.mod
	go mod tidy
	mkdir -p bin/
	go build -o ./bin/ .

.PHONY: clone clean build
