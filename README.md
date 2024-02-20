# ADSB1090

ADSB 1090 is a Mode S decoder specifically designed for RTLSDR devices.

This is a fork from the project https://github.com/antirez/dump1090 from Salvatore Sanfilippo.

The code has been cleaned to remove all net features, and has been wrapped in goland code.

This code uses:
* Leaflet: https://leafletjs.com/
* leaflet Rotated Marker: https://github.com/bbecquet/Leaflet.RotatedMarker 

## Prerequisites

* You must have a sane installation of Docker.
* Install qemu for other architecture support through docker:

```bash
sudo apt-get install -y qemu qemu-user-static
```

## Build

1. Build the drivers
2. Build the application

```bash
# Launch this command one time.
docker run --privileged --rm tonistiigi/binfmt --install all

# Build application.
make build-driver
make build-adsb1090
```

This will generate debian packages in folder `build/deb*`.

## Development

You should install Visual Studio Code (https://code.visualstudio.com/download), with extension `Dev Containers`.
Open the container and all tooling and lib will be available

## Architecture

![Diagram](archi.png)

```plantuml
@startuml
state "**Application**" as application {
    state "**source.Starter**" as source_starter {
        source_starter: Analyze encoded data and request processing

        state "File" as source_starter_file
        state "Reader" as source_starter_reader
        state "RTL28xxx" as source_starter_rtl28xxx
    }
    state "**processor.Processer**" as processor_processer {
        processor_processer: Process raw data

        state "Decoder" as processor_processer_decoder
        state "empty" as processor_processer_empty
        state "Raw" as processor_processer_raw
    }

    state "**transport.Transporter**" as transport_transporter {
        transport_transporter: Transport data

        state "File" as transport_transporter_File
        state "HTTP" as transport_transporter_http
        state "Net" as transport_transporter_net
        state "Screen" as transport_transporter_screen
    }

    state "**serialize.Serializer**" as serialize_serializer {
        serialize_serializer: Serialize data

        state "BaseStation" as serialize_serializer_basestation
        state "JSON" as serialize_serializer_json
        state "NMEA" as serialize_serializer_nmea
        state "None" as serialize_serializer_none
        state "Text" as serialize_serializer_text
    }

    state "**database.ChainedStorage**" as database_chainedstorage {
        database_chainedstorage: store ordered elements
        database_chainedstorage: //=> store squitter messages//
    }
    state "**database.ElementStorage**" as database_elementstorage {
        database_elementstorage: store any elements
        database_chainedstorage: //=> store all AC (from external source)//
    }

    state "model.Squitter" as model_squitter {
        state "ExtendedSquitter" as extended_squitter {
            extended_squitter: * AirbornePosition
            extended_squitter: * AirborneVelocity
            extended_squitter: * ExtendedSquitter
            extended_squitter: * Identification
            extended_squitter: * OperationStatus
            extended_squitter: * SurfacePosition
            extended_squitter: * Aircraft
            extended_squitter: * Position
        }

        state "ShortSquitter" as short_squitter {

        }
    }
}

source_starter --> processor_processer
processor_processer_decoder --> transport_transporter
transport_transporter_File --> serialize_serializer
transport_transporter_http --> serialize_serializer
transport_transporter_net --> serialize_serializer
transport_transporter_screen --> serialize_serializer
source_starter_file -> source_starter_reader
processor_processer_decoder --> database_chainedstorage
processor_processer_decoder --> model_squitter
transport_transporter_http --> database_elementstorage

@enduml
```