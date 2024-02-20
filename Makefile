.PHONY: build-adsb1090 build-driver clean

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

clean:
	$(MAKE) -C build clean