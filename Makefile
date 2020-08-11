.PHONY: build buildarm clean
VERSION := $(shell git describe --always | sed -e "s/^v//")
TARGETARCH=

dps-client:
	@echo "Compiling source"
	@mkdir -p build
	$(TARGETARCH) go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/dps-client cmd/dps-client/main.go

armv5: TARGETARCH=env GOOS=linux GOARCH=arm GOARM=5
armv5: dps-client

armv7: TARGETARCH=env GOOS=linux GOARCH=arm GOARM=7
armv7: dps-client

mips: TARGETARCH=env GOOS=linux GOARCH=mips
mips: dps-client

mipsle: TARGETARCH=env GOOS=linux GOARCH=mipsle
mipsle: dps-client

#Use upx to pack the binary to decrease the size
smaller:
	upx build/dps-client

clean:
	@echo "Cleaning up workspace"
	@rm -rf build
