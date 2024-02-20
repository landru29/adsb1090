// Package application is the main application.
package application

import (
	"context"
	"log/slog"

	"github.com/landru29/adsb1090/internal/config"
	"github.com/landru29/adsb1090/internal/input"
	"github.com/landru29/adsb1090/internal/input/implementations"
	"github.com/landru29/adsb1090/internal/processor"
)

// App is the main application.
type App struct {
	starter    input.Starter
	log        *slog.Logger
	processors []processor.Processer
}

// New creates a new application.
func New(
	log *slog.Logger,
	cfg *config.Config,
	processors []processor.Processer,
) (*App, error) {
	log.Info("Initializing application")

	implementations.InitTables()

	output := &App{
		log:        log,
		processors: processors,
	}

	switch {
	case cfg.FixturesFilename != "":
		// Source is a file
		opts := []implementations.FileConfigurator{}
		if cfg.FixtureLoop {
			opts = append(opts, implementations.WithLoop())
		}

		output.starter = implementations.NewFile(cfg.FixturesFilename, opts...)

		return output, nil
	default:
		opts := []implementations.RTL28Configurator{}

		if cfg.DeviceIndex > 0 {
			opts = append(opts, implementations.WithDeviceIndex(int(cfg.DeviceIndex)))
		}

		if cfg.EnableAGC {
			opts = append(opts, implementations.WithAGC())
		}

		if cfg.Frequency > 0 {
			opts = append(opts, implementations.WithFrequency(cfg.Frequency))
		}

		if cfg.Gain > 0 {
			opts = append(opts, implementations.WithGain(cfg.Gain))
		}

		output.starter = implementations.New(opts...)

		return output, nil
	}
}

// Start is the application entrypoint.
func (a *App) Start(ctx context.Context) error {
	a.log.Info("Starting application")

	return a.starter.Start(ctx, a.processors...)
}
