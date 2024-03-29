package implementations

/*
  #cgo LDFLAGS: -lrtlsdr -lm
  #include "rtlsdr.h"
  #include <malloc.h>

*/
import "C"

import (
	"context"
	"fmt"
	"unsafe"

	localcontext "github.com/landru29/adsb1090/internal/input/context"
	"github.com/landru29/adsb1090/internal/logger"
	"github.com/landru29/adsb1090/internal/processor"
)

var debug string //nolint: gochecknoglobals

// Device is a RTL-SDR device.
type Device struct {
	dev        *C.rtlsdr_dev_t
	processors []processor.Processer
}

// InitTables generates tables for data extract.
func InitTables() {
	C.initTables(map[string]C.int{"false": 0, "true": 1}[debug])
}

// DeviceCount searches for a compatible device.
func DeviceCount() uint32 {
	return uint32(C.rtlsdr_get_device_count())
}

// DeviceUsbStrings gets USB device strings.
func DeviceUsbStrings(index uint32) (string, string, string, error) {
	serial := (*C.char)(C.malloc(256))   //nolint: gomnd
	manufact := (*C.char)(C.malloc(256)) //nolint: gomnd
	product := (*C.char)(C.malloc(256))  //nolint: gomnd

	defer func() {
		C.free(unsafe.Pointer(serial))
		C.free(unsafe.Pointer(manufact))
		C.free(unsafe.Pointer(product))
	}()

	if intErr := C.rtlsdr_get_device_usb_strings(C.uint32_t(index), manufact, product, serial); intErr != 0 {
		return "", "", "", fmt.Errorf("RtlsdrGetDeviceUsbStrings: %d", intErr)
	}

	return C.GoString(manufact), C.GoString(product), C.GoString(serial), nil
}

// OpenDevice opens the device.
func OpenDevice(index uint32, processors []processor.Processer) (*Device, error) {
	output := Device{
		processors: processors,
	}

	if intErr := C.rtlsdr_open(&output.dev, C.uint32_t(index)); intErr != 0 { //nolint: gocritic,nlreturn
		return nil, fmt.Errorf("RtlsdrOpen: %d", intErr)
	}

	return &output, nil
}

// Close closes the device.
func (d *Device) Close() error {
	if intErr := C.rtlsdr_close(d.dev); intErr != 0 { //nolint: nlreturn
		return fmt.Errorf("RtlsdrClose: %d", intErr)
	}

	return nil
}

// SetTunerGainMode sets the gain mode (automatic/manual) for the device.
// Manual gain mode must be enabled for the gain setter function to work.
func (d *Device) SetTunerGainMode(manual bool) error {
	if intErr := C.rtlsdr_set_tuner_gain_mode(
		d.dev,
		map[bool]C.int{ //nolint: nlreturn
			true:  1,
			false: 0,
		}[manual],
	); intErr != 0 {
		return fmt.Errorf("RtlsdrSetTunerGainMode: %d", intErr)
	}

	return nil
}

// TunerGains gets a list of gains supported by the tuner.
//
// NOTE: The gains argument must be preallocated by the caller. If NULL is
// being given instead, the number of available gain values will be returned.
func (d *Device) TunerGains() ([]int, error) {
	gains := (*C.int)(C.malloc(100 * C.sizeof_int)) //nolint: gomnd

	size := C.rtlsdr_get_tuner_gains(d.dev, gains) //nolint: nlreturn
	if size < 0 {
		return nil, fmt.Errorf("RtlsdrGetTunerGains: %d", size)
	}

	castedGains := (*[100]C.int)(unsafe.Pointer(gains))

	outGains := make([]int, 100) //nolint: gomnd
	for idx := range outGains {
		outGains[idx] = int(castedGains[idx])
	}

	return outGains[:size], nil
}

// SetTunerGain sets the gain for the device.
// Manual gain mode must be enabled for this to work.
//
// Valid gain values (in tenths of a dB) for the E4000 tuner:
// -10, 15, 40, 65, 90, 115, 140, 165, 190, 215, 240, 290, 340, 420, 430, 450, 470, 490
func (d *Device) SetTunerGain(gain float64) error {
	if intErr := C.rtlsdr_set_tuner_gain(d.dev, C.int(gain*10)); intErr != 0 { //nolint: gomnd,nlreturn
		return fmt.Errorf("RtlsdrSetTunerGain: %d", intErr)
	}

	return nil
}

// TunerGain read the current gain on the tuner.
func (d *Device) TunerGain() float64 {
	return float64(C.rtlsdr_get_tuner_gain(d.dev) / 10.0)
}

// SetFreqCorrection sets the frequency correction value for the device.
func (d *Device) SetFreqCorrection(partsPerMillion int) error {
	if intErr := C.rtlsdr_set_freq_correction(d.dev, C.int(partsPerMillion)); intErr != 0 { //nolint: gomnd,nlreturn,nolintlint,lll
		return fmt.Errorf("RtlsdrSetFreqCorrection: %d", intErr)
	}

	return nil
}

// SetAgcMode enables or disables the internal digital AGC of the RTL2832.
func (d *Device) SetAgcMode(on bool) error {
	if intErr := C.rtlsdr_set_agc_mode(
		d.dev,
		map[bool]C.int{ //nolint: nlreturn
			true:  1,
			false: 0,
		}[on],
	); intErr != 0 {
		return fmt.Errorf("RtlsdrSetAgcMode: %d", intErr)
	}

	return nil
}

// SetCenterFreq ...
func (d *Device) SetCenterFreq(freq uint32) error {
	if intErr := C.rtlsdr_set_center_freq(d.dev, C.uint32_t(freq)); intErr != 0 { //nolint: gomnd,nlreturn,nolintlint,lll
		return fmt.Errorf("RtlsdrSetCenterFreq: %d", intErr)
	}

	return nil
}

// SetSampleRate sets the sample rate for the device, also selects the baseband filters
// according to the requested sample rate for tuners where this is possible.
//
// possible values of rate are:
//
//	225001 - 300000 Hz
//	900001 - 3200000 Hz
//	sample loss is to be expected for rates > 2400000
func (d *Device) SetSampleRate(rate uint32) error {
	if intErr := C.rtlsdr_set_sample_rate(d.dev, C.uint32_t(rate)); intErr != 0 { //nolint: gomnd,nlreturn,nolintlint,lll
		return fmt.Errorf("RtlsdrSetSampleRate: %d", intErr)
	}

	return nil
}

// ResetBuffer resets the device internal buffer.
func (d *Device) ResetBuffer() error {
	if intErr := C.rtlsdr_reset_buffer(d.dev); intErr != 0 { //nolint: nlreturn
		return fmt.Errorf("RtlsdrResetBuffer: %d", intErr)
	}

	return nil
}

// ReadAsync reads samples from the device asynchronously. This function will block until
// it is being canceled using rtlsdr_cancel_async()
func (d *Device) ReadAsync(ctx context.Context, bufNum uint32, bufLen uint32) error {
	log, logFound := logger.Logger(ctx)

	if logFound {
		log.Info("Launching an asynchronous read on the device")
	}

	if intErr := C.rtlsdrReadAsync(
		d.dev,
		localcontext.New(ctx, d.processors).Ccontext,
		C.uint32_t(bufNum),
		C.uint32_t(bufLen), //nolint: nlreturn
	); intErr != 0 {
		return fmt.Errorf("RtlsdrReadAsync: %d", intErr)
	}

	if logFound {
		log.Info("asynchronous read ended")
	}

	return nil
}

func processRaw(data []byte, cContext unsafe.Pointer) {
	cstr := (*C.uchar)(unsafe.Pointer(C.CString(string(data))))

	C.rtlsdrProcessRaw(cstr, C.uint(len(data)), cContext)

	C.free(unsafe.Pointer(cstr))
}

//export goRtlsrdData
func goRtlsrdData(buf *C.uchar, length C.uint32_t, cCtx *C.void) C.int {
	ctx := localcontext.FromPtr(unsafe.Pointer(cCtx))
	processors := localcontext.Processor(ctx)

	mySlice := C.GoBytes(unsafe.Pointer(buf), C.int(length)) //nolint: nlreturn

	for _, processor := range processors {
		if err := processor.Process(mySlice); err != nil {
			return -1
		}
	}

	return 0
}
