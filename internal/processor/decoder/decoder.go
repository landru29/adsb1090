// Package decoder is the default processor.
package decoder

import (
	"context"
	"log/slog"
	"time"

	"github.com/landru29/adsb1090/internal/aircraftdb"
	"github.com/landru29/adsb1090/internal/database"
	"github.com/landru29/adsb1090/internal/errors"
	"github.com/landru29/adsb1090/internal/model"
	"github.com/landru29/adsb1090/internal/transport"
)

const (
	// ErrReferenceAircraftNotFound is when the reference aircraft was not found in the world database.
	ErrReferenceAircraftNotFound errors.Error = "aircraft not found in the world database"
)

// Configurator is the Process configurator.
type Configurator func(*Process)

// Process is the data processor.
type Process struct {
	ExtendedSquitters     *database.ChainedStorage[model.ICAOAddr, model.QualifiedMessage]
	log                   *slog.Logger
	dbLifeTime            time.Duration
	transporters          []transport.Transporter
	aircraftWorldDatabase aircraftdb.Database
}

// New creates a data processor.
func New(ctx context.Context, log *slog.Logger, opts ...Configurator) *Process {
	process := &Process{
		log:          log,
		transporters: []transport.Transporter{},
	}

	for _, opt := range opts {
		opt(process)
	}

	process.ExtendedSquitters = database.NewChainedStorage[model.ICAOAddr, model.QualifiedMessage](
		ctx,
		database.ChainedWithLifetime[model.ICAOAddr, model.QualifiedMessage](process.dbLifeTime),
	)

	return process
}

// WithAircraftWorldDatabase sets the aircraft world database.
func WithAircraftWorldDatabase(aircraftWorldDatabase aircraftdb.Database) Configurator {
	return func(process *Process) {
		process.aircraftWorldDatabase = aircraftWorldDatabase
	}
}

// WithDatabaseLifetime sets the lifetime of database elements.
func WithDatabaseLifetime(dbLifeTime time.Duration) Configurator {
	return func(process *Process) {
		process.dbLifeTime = dbLifeTime
	}
}

// WithTransporter add a new transporter.
func WithTransporter(transporter transport.Transporter) Configurator {
	return func(process *Process) {
		process.transporters = append(process.transporters, transporter)
	}
}

// Process implements source.Processor the interface.
func (p Process) Process(data []byte) error {
	modes := model.ModeS(data)

	log := p.log.With("message", modes.String())

	squitter, err := modes.QualifiedMessage()
	if err != nil {
		return err
	}

	if _, isExtended := squitter.(model.LongMessage); isExtended {
		if err := modes.CheckSum(); err != nil {
			return err
		}
	}

	log.Info("processing message")

	icaoAddress := squitter.AircraftAddress()

	aircraftReference, found := p.aircraftWorldDatabase[icaoAddress]
	if !found {
		return ErrReferenceAircraftNotFound
	}

	log = log.
		With("registration", aircraftReference.Registration).
		With("model", aircraftReference.Model).
		With("operator", aircraftReference.Operator).
		With("manufacturer", aircraftReference.ManufacturerName)

	log.Info("aircraft found")

	p.ExtendedSquitters.Add(icaoAddress, squitter)

	aircraft := buildAircraft(log, p.ExtendedSquitters.Elements(icaoAddress), aircraftReference)

	for _, transporter := range p.transporters {
		if err := transporter.Transport(aircraft); err != nil {
			log.Error("transport", "msg", err)
		}
	}

	return nil
}
