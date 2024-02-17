// Package main is the main application.
package main

import (
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/landru29/adsb1090/internal/aircraftdb"
	"github.com/landru29/adsb1090/internal/application"
	"github.com/landru29/adsb1090/internal/database"
	"github.com/landru29/adsb1090/internal/logger"
	"github.com/landru29/adsb1090/internal/model"
	"github.com/landru29/adsb1090/internal/processor"
	"github.com/landru29/adsb1090/internal/processor/decoder"
	"github.com/landru29/adsb1090/internal/serialize"
	"github.com/landru29/adsb1090/internal/serialize/nmea"
	"github.com/landru29/adsb1090/internal/transport/net"
	"github.com/spf13/cobra"
)

const (
	defaultNMEAmid                        = 226
	defaultFrequency                      = 1090000000
	defaultDatabaseLifetime time.Duration = time.Minute
)

func rootCommand() *cobra.Command { //nolint: funlen
	var (
		app                    *application.App
		config                 application.Config
		httpConf               httpConfig
		transportScreen        string
		transportFile          string
		nmeaMid                uint16
		nmeaVessel             vessel = nmea.VesselTypeAircraft
		availableSerializers   []serialize.Serializer
		loop                   bool
		settings               *application.Settings
		refAircraftDatabaseURL string
	)

	udpConf := net.NewProtocol("udp")
	tcpConf := net.NewProtocol("tcp")

	rootCommand := &cobra.Command{
		Use:   "adsb1090",
		Short: "adsb1090",
		Long:  "adsb1090 main command",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			settings, err = application.UserSettings()
			if err != nil {
				return err
			}

			log := slog.New(slog.NewTextHandler(cmd.OutOrStdout(), nil))

			ctx := settings.WithSettings(
				logger.WithLogger(
					cmd.Context(),
					log,
				),
			)

			cmd.SetContext(ctx)

			aircraftDB := database.NewElementStorage[model.ICAOAddr, model.Aircraft](
				ctx,
				database.ElementWithLifetime[model.ICAOAddr, model.Aircraft](config.DatabaseLifetime),
				database.ElementWithCleanCycle[model.ICAOAddr, model.Aircraft](config.DatabaseLifetime),
			)

			var serializers map[string]serialize.Serializer

			serializers, availableSerializers = provideSerializers(log, nmea.VesselType(nmeaVessel), nmeaMid)

			transporters, err := provideTransporters(
				ctx,
				log,
				availableSerializers,
				serializers,
				aircraftDB,
				httpConf,
				udpConf,
				tcpConf,
				transportScreen,
				transportFile,
			)
			if err != nil {
				return err
			}

			decoderCfg := []decoder.Configurator{
				decoder.WithDatabaseLifetime(config.DatabaseLifetime),
			}
			for _, transporter := range transporters {
				log.Info("loading transporter", "name", transporter.String())

				decoderCfg = append(decoderCfg, decoder.WithTransporter(transporter))
			}

			log.Info("loading aircraft database", "from", settings.AircraftDatabaseFile())
			aircraftWorldDatabase := aircraftdb.Database{}

			for retry := 0; retry < 2; retry++ {
				if err := aircraftWorldDatabase.Load(settings.AircraftDatabaseFile(), io.Discard); err != nil {
					if err := download(settings.AircraftDatabaseFile(), refAircraftDatabaseURL, cmd.ErrOrStderr()); err != nil {
						log.Error(fmt.Sprintf("please download the aircraft database first '%s aircraft download'", cmd.CommandPath()))

						return err
					}
				}
			}

			log.Info("found", "count", len(aircraftWorldDatabase))
			decoderCfg = append(decoderCfg, decoder.WithAircraftWorldDatabase(aircraftWorldDatabase))

			app, err = application.New(
				log,
				&config,
				[]processor.Processer{
					decoder.New(
						ctx,
						log,
						decoderCfg...,
					),
					// raw.New(log),
				},
			)

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := app.Start(ctx); err != nil {
				return err
			}

			<-ctx.Done()

			return nil
		},
	}

	rootCommand.Flags().StringVarP(
		&refAircraftDatabaseURL,
		"url",
		"u",
		"https://opensky-network.org/datasets/metadata/",
		"URL to aircraft database",
	)

	rootCommand.Flags().StringVarP(
		&config.FixturesFilename,
		"fixture-file",
		"",
		"",
		"Filename of the fixture data file",
	)

	rootCommand.Flags().Uint32VarP(
		&config.DeviceIndex,
		"device",
		"d",
		0,
		"Device index",
	)

	rootCommand.Flags().BoolVarP(
		&config.EnableAGC,
		"enable-agc",
		"a",
		false,
		"Enable AGC",
	)

	rootCommand.Flags().Uint32VarP(
		&config.Frequency,
		"frequency",
		"f",
		defaultFrequency,
		"frequency in Hz",
	)

	rootCommand.Flags().DurationVarP(
		&config.DatabaseLifetime,
		"db-lifetime",
		"",
		defaultDatabaseLifetime,
		"lifetime of elements in the AC database",
	)

	rootCommand.Flags().Float64VarP(
		&config.Gain,
		"gain",
		"g",
		0,
		"gain valid values are: 1.5, 4, 6.5, 9, 11.5, 14, 16.5, 19, 21.5, 24, 29, 34, 42, 43, 45, 47, 49",
	)

	rootCommand.Flags().VarP(
		&udpConf,
		"udp",
		"",
		"transmit data over udp (syntax: 'direction>format@host:port'; ie: --udp dial>json@192.168.1.10:8000)",
	)

	rootCommand.Flags().VarP(
		&tcpConf,
		"tcp",
		"",
		"transmit data over tcp (syntax: 'direction>format@host:port'; ie: --tcp bind>json@192.168.1.10:8000)",
	)

	rootCommand.Flags().VarP(
		&httpConf,
		"http",
		"",
		"transmit data over http (syntax: 'host:port/path'; ie: --http 0.0.0.0:8080/api)",
	)

	rootCommand.Flags().StringVarP(
		&transportScreen,
		"screen",
		"",
		"",
		"format to display output on the screen (json|nmea|text|none)",
	)

	rootCommand.Flags().VarP(
		&nmeaVessel,
		"nmea-vessel",
		"",
		"MMSI vessel (aircraft|helicopter)",
	)

	rootCommand.Flags().Uint16VarP(
		&nmeaMid,
		"nmea-mid",
		"",
		defaultNMEAmid,
		"MID (command 'mid' to list)",
	)

	rootCommand.Flags().BoolVarP(
		&loop,
		"loop",
		"",
		false,
		"With --fixture-file, read the same file in a loop",
	)

	rootCommand.Flags().StringVarP(
		&transportFile,
		"out-file",
		"",
		"",
		"format to display output on a file; ie --out-file nmea@/tmp/foo.txt",
	)

	rootCommand.AddCommand(
		aircraftCommand(),
		serializerCommand(&availableSerializers),
	)

	return rootCommand
}
