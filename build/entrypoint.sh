#!/bin/bash

set -e

for arg in "$@"
do
    case $arg in
        "build-driver" )
            cd /application
            make rtlsdr
            cp *.deb /deb/.
           exit 0
           ;;
        "build-adsb1090" )
            apt install -y /deb/*rtl*.deb

            cd /application

            make package

            echo "Packaging success !"
            exit 0
           ;;
   esac
done



