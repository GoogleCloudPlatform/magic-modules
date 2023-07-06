localsetup:
	sed -i bak -E "s~google/provider/(.*)/([0-9A-Za-z-]*)~google/provider/\1/$(PROVIDER_VERSION)~" comparison.go

ifneq ($(REPLACE_OLD),)
	go mod edit -replace google/provider/old=$(realpath ${REPLACE_OLD})
endif
ifneq ($(REPLACE_NEW),)
	go mod edit -replace google/provider/new=$(realpath ${REPLACE_NEW})
endif