ARCH := $(shell dpkg --print-architecture)
CURRENT_FOLDER := $(shell pwd)
DEB_DEST ?= "$(PWD)/deb"

.PHONY: build-adsb1090 build-driver clean build package rtlsdr deps

package: adsb1090-$(ARCH).deb

build: deps
	mkdir -p $(DEB_DEST)/adsb1090/usr/local/bin && \
	go build -buildvcs=false -o $(DEB_DEST)/adsb1090/usr/local/bin ./cmd/adsb1090/...

deps: 
	go mod download

adsb1090-$(ARCH).deb: build $(DEB_DEST)
	cp -r build/DEBIAN $(DEB_DEST)/adsb1090/. && \
	cp -r build/etc $(DEB_DEST)/adsb1090/. && \
	cat build/DEBIAN/control | sed -e "s/xxxxx/${ARCH}/g">$(DEB_DEST)/adsb1090/DEBIAN/control && \
	cd $(DEB_DEST) && \
	dpkg-deb --build --root-owner-group adsb1090 && \
	rm -rf $(DEB_DEST)/adsb1090

rtlsdr: $(DEB_DEST)
	git clone --depth 1 --branch v2.0.1 https://github.com/osmocom/rtl-sdr.git  /tmp/rtl-sdr && \
	cd /tmp/rtl-sdr && \
	dpkg-buildpackage -us -uc && \
	cp /tmp/*rtl*.deb $(DEB_DEST) && \
	rm -rf /tmp/rtl-sdr

$(DEB_DEST): 
	mkdir -p $(DEB_DEST)

lint:
	golangci-lint run ./...

test:
	gotest ./...

gen:
	go generate ./...

build-adsb1090: build-adsb1090.armv6l build-adsb1090.armhf build-adsb1090.amd64 build-adsb1090.arm64

build-driver: build-driver.armv6l build-driver.armhf build-driver.amd64 build-driver.arm64 

build-driver.%:
	$(MAKE) -C build $@

build-adsb1090.%:
	$(MAKE) -C build $@

debug.%:
	QEMU_CPU=$(QEMU_CPU) $(MAKE) -C build $@

clean:
	$(MAKE) -C build clean
	rm -rf build/usr
	rm -f *.deb