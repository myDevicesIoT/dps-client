.PHONY: dps-client armv5 armv7 mips mipsle smaller multitech tektelic gemtek clean
VERSION := $(shell git describe --always | sed -e "s/^v//")
ARCH := $(shell uname -m)
TARGETARCH=
LDFLAGS="-s -w -X main.version=$(VERSION)"

dps-client:
	@echo "Compiling source"
	@mkdir -p build
	$(TARGETARCH) go build $(GO_EXTRA_BUILD_ARGS) -ldflags $(LDFLAGS) -o build/$(ARCH)/dps-client cmd/dps-client/main.go

armv5: TARGETARCH=env GOOS=linux GOARCH=arm GOARM=5
armv5: ARCH=armv5
armv5: dps-client

armv7: TARGETARCH=env GOOS=linux GOARCH=arm GOARM=7 ARCH=armv7
armv7: ARCH=armv7
armv7: dps-client

mips: TARGETARCH=env GOOS=linux GOARCH=mips ARCH=mips
mips: ARCH=mips
mips: dps-client

mipsle: TARGETARCH=env GOOS=linux GOARCH=mipsle ARCH=mipsle
mipsle: ARCH=mipsle
mipsle: dps-client

#Use upx to pack the binary to decrease the size
smaller:
	upx build/*/dps-client

multitech: armv5
	upx build/armv5/dps-client
	cd packaging/multitech; ./package.sh $(VERSION)

tektelic: armv5
	upx build/armv5/dps-client
	cd packaging/tektelic; ./package.sh $(VERSION)

gemtek: mipsle
	upx build/mipsle/dps-client
	cd packaging/gemtek; ./package.sh $(VERSION)	

clean:
	@echo "Cleaning up workspace"
	@rm -rf build
