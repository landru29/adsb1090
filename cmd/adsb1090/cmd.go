// Package main is the main application.
package main

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/landru29/adsb1090/internal/aircraftdb"
	"github.com/landru29/adsb1090/internal/application"
	conf "github.com/landru29/adsb1090/internal/config"
	"github.com/landru29/adsb1090/internal/database"
	"github.com/landru29/adsb1090/internal/logger"
	"github.com/landru29/adsb1090/internal/model"
	"github.com/landru29/adsb1090/internal/processor"
	"github.com/landru29/adsb1090/internal/processor/decoder"
	"github.com/landru29/adsb1090/internal/serialize"
	"github.com/landru29/adsb1090/internal/serialize/nmea"
	"github.com/spf13/cobra"
)

func rootCommand() (*cobra.Command, error) { //nolint: funlen
	var (
		app                  *application.App
		availableSerializers []serialize.Serializer
		config               *conf.Config
	)

	rootCommand := &cobra.Command{
		Use:   "adsb1090",
		Short: "adsb1090",
		Long:  "adsb1090 main command",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			log := slog.New(slog.NewTextHandler(cmd.OutOrStdout(), nil))

			ctx := logger.WithLogger(
				cmd.Context(),
				log,
			)

			cmd.SetContext(ctx)

			aircraftDB := database.NewElementStorage[model.ICAOAddr, model.Aircraft](
				ctx,
				database.ElementWithLifetime[model.ICAOAddr, model.Aircraft](config.DatabaseLifetime),
				database.ElementWithCleanCycle[model.ICAOAddr, model.Aircraft](config.DatabaseLifetime),
			)

			var serializers map[string]serialize.Serializer

			serializers, availableSerializers = provideSerializers(log, nmea.VesselType(config.NmeaVessel), config.NmeaMid)

			transporters, err := provideTransporters(
				ctx,
				log,
				availableSerializers,
				serializers,
				aircraftDB,
				config.HTTPConf,
				config.UDPConf,
				config.TCPConf,
				config.TransportScreen,
				config.TransportFile,
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

			log.Info("loading aircraft database", "from", config.AircraftDatabaseFile())
			aircraftWorldDatabase := aircraftdb.Database{}

			for retry := 0; retry < 2; retry++ {
				if err := aircraftWorldDatabase.Load(config.AircraftDatabaseFile(), io.Discard); err != nil {
					if err := download(config.AircraftDatabaseFile(), config.RefAircraftDatabaseURL, cmd.ErrOrStderr()); err != nil {
						log.Error(fmt.Sprintf("please download the aircraft database first '%s aircraft download'", cmd.CommandPath()))

						return err
					}
				}
			}

			log.Info("found", "count", len(aircraftWorldDatabase))
			decoderCfg = append(decoderCfg, decoder.WithAircraftWorldDatabase(aircraftWorldDatabase))

			app, err = application.New(
				log,
				config,
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

	var err error

	config, err = conf.UserSettings(rootCommand.Flags())
	if err != nil {
		return nil, err
	}

	rootCommand.AddCommand(
		aircraftCommand(config),
		serializerCommand(&availableSerializers),
		configCommand(config),
	)

	return rootCommand, nil
}
