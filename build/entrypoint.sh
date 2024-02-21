#!/bin/bash

set -e

for arg in "$@"
do
    case $arg in
        "build-driver" )
            mkdir -p /project/c && cd /project/c
            git clone --depth 1 --branch v2.0.1 https://github.com/osmocom/rtl-sdr.git
            cd rtl-sdr
            dpkg-buildpackage -us -uc
            cp /project/c/*.deb /deb/.
           exit 0
           ;;
        "build-adsb1090" )
            cd /deb
            apt install -y \
                ./librtlsdr*.deb \
                ./librtlsdr-dev*.deb \
                ./rtl-sdr*.deb \
                ./rtl-sdr-dbgsym*.deb 
            cd /application

            case $(arch) in
                "x86_64" )
                    ARCH=amd64
                    ;;
                "aarch64" )
                    ARCH=arm64
                    ;;
                "armv7l" )
                    ARCH=armhf
                    ;;
                "armv6l" )
                    ARCH=armv6l
                    ;;
            esac

            mkdir -p /adsb1090-${ARCH}/usr/local/bin

            cp -r /application/build/DEBIAN /adsb1090-${ARCH}/.
            cp -r /application/build/etc /adsb1090-${ARCH}/.

            cat /adsb1090-${ARCH}/DEBIAN/control | sed -e "s/amd64/${ARCH}/g">/tmp/control
            cp /tmp/control /adsb1090-${ARCH}/DEBIAN/control

            export CGO_ENABLED=1
            export GOPROXY=direct

            echo "== Getting dependencies =="

            go mod download

            echo "== Building =="

            go build -o /adsb1090-${ARCH}/usr/local/bin/adsb1090 -buildvcs=false ./cmd/adsb1090/...

            echo "== Packaging =="

            dpkg-deb --build --root-owner-group /adsb1090-${ARCH}
            cp /adsb1090-${ARCH}.deb /deb/.

            echo "Packaging success !"
            exit 0
           ;;
   esac
done



