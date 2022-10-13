localsetup:
ifneq (,$(filter $(PROVIDER_VERSION),"google" "google-beta"))
	echo "setting up environment for google-beta"
	PROVIDER_VERSION="google-beta"
else
	echo "setting up environment for $(PROVIDER_VERSION)"
endif

	sed -i bak -E "s~google/provider/(.*)/([0-9A-Za-z-]*)~google/provider/\1/$(PROVIDER_VERSION)~" comparison.go

ifneq ($(REPLACE_OLD),)
	go mod edit -replace google/provider/old=$(realpath ${REPLACE_OLD})
endif
ifneq ($(REPLACE_NEW),)
	go mod edit -replace google/provider/new=$(realpath ${REPLACE_NEW})
endif