// Package implementations is the RTL28xxx data source.
package implementations

import (
	"context"

	"github.com/landru29/adsb1090/internal/errors"
	"github.com/landru29/adsb1090/internal/logger"
	"github.com/landru29/adsb1090/internal/processor"
)

const (
	// ErrNoDeviceFound is when no device is found.
	ErrNoDeviceFound errors.Error = "no device found"

	modeSfrequency = 1090000000
	sampleRate     = 2000000

	asyncBufNumber = 12
	dataLen        = (16 * 32 * 512) /* 256k */ //nolint: gomnd
)

// RTL28Configurator is the Source configurator.
type RTL28Configurator func(*RTL28xxx)

// RTL28xxx is the data source process.
type RTL28xxx struct {
	deviceIndex uint32
	frequency   uint32
	gain        float64
	enableAGC   bool

	dev *Device
}

// New creates a new data source process.
func New(opts ...RTL28Configurator) *RTL28xxx {
	output := &RTL28xxx{
		deviceIndex: 0,
		frequency:   modeSfrequency,
		gain:        0,
		enableAGC:   false,
	}

	for _, opt := range opts {
		opt(output)
	}

	return output
}

// Start implements the input.Starter interface.
func (s *RTL28xxx) Start(ctx context.Context, processors ...processor.Processer) error { //nolint: cyclop
	log, loggerFound := logger.Logger(ctx)

	deviceCount := DeviceCount()
	if deviceCount == 0 {
		return ErrNoDeviceFound
	}

	deviceIndex := uint32(0)

	if s.deviceIndex < deviceCount {
		deviceIndex = s.deviceIndex
	}

	device, err := OpenDevice(deviceIndex, processors)
	if err != nil {
		return err
	}

	if loggerFound {
		log.Info("device found")
	}

	s.dev = device

	if err := s.dev.SetCenterFreq(modeSfrequency); err != nil {
		return err
	}

	if loggerFound {
		log.Info("configuring device", "frequency", modeSfrequency)
	}

	if err := s.dev.SetSampleRate(sampleRate); err != nil {
		return err
	}

	if loggerFound {
		log.Info("configuring sample rate", "rate", sampleRate)
	}

	if err := s.dev.SetAgcMode(s.enableAGC); err != nil {
		return err
	}

	if loggerFound {
		log.Info("configuring AGC", "agc", s.enableAGC)
	}

	if err := s.dev.SetTunerGainMode(s.gain > 0); err != nil {
		return err
	}

	if loggerFound {
		log.Info("configuring gain mode", "mode", map[bool]string{false: "auto", true: "manual"}[s.gain > 0])
	}

	if s.gain > 0 {
		if err := s.dev.SetTunerGain(s.gain); err != nil {
			return err
		}

		if loggerFound {
			log.Info("configuring gain", "gain", s.gain)
		}
	}

	if err := s.dev.ResetBuffer(); err != nil {
		return err
	}

	if loggerFound {
		log.Info("device ready", "gain", s.dev.TunerGain())
	}

	return s.dev.ReadAsync(ctx, asyncBufNumber, dataLen)
}

// WithDeviceIndex configures the device index.
func WithDeviceIndex(index int) RTL28Configurator {
	return func(s *RTL28xxx) {
		if index > 0 {
			s.deviceIndex = uint32(index)
		}
	}
}

// WithFrequency configures the frequency.
func WithFrequency(frequency uint32) RTL28Configurator {
	return func(s *RTL28xxx) {
		s.frequency = frequency
	}
}

// WithGain configures the gain.
func WithGain(gain float64) RTL28Configurator {
	return func(s *RTL28xxx) {
		s.gain = gain
	}
}

// WithAGC enables AGC.
func WithAGC() RTL28Configurator {
	return func(s *RTL28xxx) {
		s.enableAGC = true
	}
}
