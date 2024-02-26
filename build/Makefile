image_name=driver-rtl-sdr
PWD:=$(shell pwd)
ifneq ($(QEMU_CPU),)
	ENV=--env QEMU_CPU=$(QEMU_CPU)
endif

.PHONY: clean

build-adsb1090.%: TARGET_NAME=$(subst build-adsb1090.,,$@)
build-adsb1090.%: docker.% deb.%
	docker run \
		--env DEB_DEST=/deb \
		--env CGO_ENABLED=1 \
		--env GOPROXY=direct \
	    --mount type=bind,source=$(PWD)/deb.$(TARGET_NAME),target=/deb \
		--volume $(PWD)/..:/application:ro \
		$(image_name)-$(TARGET_NAME) \
		build-adsb1090 && \
		find . -type d -iname 'deb.*' -exec tar czf {}.tgz {} \;

build-driver.%: TARGET_NAME=$(subst build-driver.,,$@)
build-driver.%: docker.% deb.%
	docker run \
		--env DEB_DEST=/deb \
		--env CGO_ENABLED=1 \
		--env GOPROXY=direct \
	    --mount type=bind,source=$(PWD)/deb.$(TARGET_NAME),target=/deb \
		$(image_name)-$(TARGET_NAME) \
		build-driver && \
		find . -type d -iname 'deb.*' -exec tar czf {}.tgz {} \;


debug.%: TARGET_NAME=$(subst debug.,,$@)
debug.%: docker.% deb.%
	docker run -it \
	    $(ENV) \
		--env DEB_DEST=/deb \
		--env CGO_ENABLED=1 \
		--env GOPROXY=direct \
	    --mount type=bind,source=$(PWD)/deb.$(TARGET_NAME),target=/deb \
		--volume $(PWD)/..:/application:ro \
		--entrypoint bash \
		$(image_name)-$(TARGET_NAME)

docker.%:
	docker build -t $(image_name)-$(subst docker.,,$@) -f $(subst docker,Dockerfile,$@) .

clean:
	docker image rm -f $(image_name)-arm64
	docker image rm -f $(image_name)-amd64
	rm -rf deb.*

deb.%:
	mkdir -p $@

