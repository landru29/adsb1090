package model

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/landru29/adsb1090/internal/binary"
	localerrors "github.com/landru29/adsb1090/internal/errors"
)

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
// ┃                                  Mode S                                    ┃
// ┠┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
// ┃                                    112                                     ┃
// ┣━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┫
// ┃ DF  |                          squitter                           | Parity ┃
// ┠┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┨
// ┃  5  |                             27 / 83                         |   24   ┃
// ┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┛
//
// DF = 0, 4, 5, 11
// ┏━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┓
// ┃ DF  |                        short squitter                       | Parity ┃
// ┠┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┨
// ┃  5  |                             27                              |   24   ┃
// ┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┛
//
// DF = 16-21, 24
// ┏━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┓
// ┃ DF  |                      Extended squitter                      | Parity ┃
// ┠┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┨
// ┃  5  |                             83                              |   24   ┃
// ┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┛

//
// ┏━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━┓
// ┃ name | Description                 | bits ┃
// ┣━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┿━━━━━━┫
// ┃ DF   | Downlink Format             |   5  ┃
// ┗━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━┛
//
//

// DownlinkFormat is the 5 first bits of an ADSB message.
type DownlinkFormat uint8

const (
	shortSquitterBitLength    = 56
	extendedSquitterBitLength = 112

	// DownlinkFormatShortAirAirSurveillance is Short air-air surveillance (ACAS) => message size: 56 bits.
	DownlinkFormatShortAirAirSurveillance DownlinkFormat = 0
	// DownlinkFormatAltitudeReply is Altitude reply => message size: 56 bits.
	DownlinkFormatAltitudeReply DownlinkFormat = 4
	// DownlinkFormatIdentityReply is Identity reply => message size: 56 bits.
	DownlinkFormatIdentityReply DownlinkFormat = 5
	// DownlinkFormatAllCallReply is All-call reply => message size: 56 bits.
	DownlinkFormatAllCallReply DownlinkFormat = 11
	// DownlinkFormatLongAirAirSurveillance is Long air-air surveillance (ACAS) => message size: 112 bits.
	DownlinkFormatLongAirAirSurveillance DownlinkFormat = 16
	// DownlinkFormatExtendedSquitter is Extended squitter => message size: 112 bits.
	DownlinkFormatExtendedSquitter DownlinkFormat = 17
	// DownlinkFormatExtendedSquitterNonTransponder is Extended squitter, non transponder => message size: 112 bits.
	DownlinkFormatExtendedSquitterNonTransponder DownlinkFormat = 18
	// DownlinkFormatMilitaryExtendedSquitter is Military extended squitter => message size: 112 bits.
	DownlinkFormatMilitaryExtendedSquitter DownlinkFormat = 19
	// DownlinkFormatCommBWithAltitudeReply is Comm-B, with altitude reply => message size: 112 bits.
	DownlinkFormatCommBWithAltitudeReply DownlinkFormat = 20
	// DownlinkFormatCommBWithIdentityReply is Comm-B, with identity reply => message size: 112 bits.
	DownlinkFormatCommBWithIdentityReply DownlinkFormat = 21
	// DownlinkFormatCommDExtendedLengthMessage is Comm-D, extended length message => message size: 112 bits.
	DownlinkFormatCommDExtendedLengthMessage DownlinkFormat = 24

	// ErrUnsupportedFormat is when the mode-s format is not supported.
	ErrUnsupportedFormat localerrors.Error = "unsupported  format"
	// ErrWrongCRC is when a wrong CRC was encountered.
	ErrWrongCRC localerrors.Error = "wrong CRC"
)

// ModeS is a ModeS frame.
type ModeS []byte

// DownlinkFormat is the DF.
func (m ModeS) DownlinkFormat() DownlinkFormat {
	return DownlinkFormat((m[0] & 0xf8) >> 3) //nolint: gomnd
}

// ParityInterrogator is the Parity.
func (m ModeS) ParityInterrogator() uint32 {
	length := len(m)

	return (uint32(m[length-3]) << 16) + (uint32(m[length-2]) << 8) + //nolint: gomnd
		uint32(m[length-1])
}

// QualifiedMessage is the qualified message.
func (m ModeS) QualifiedMessage() (QualifiedMessage, error) { //nolint: ireturn
	downlinkFormat := m.DownlinkFormat()

	if (downlinkFormat == DownlinkFormatExtendedSquitter ||
		downlinkFormat == DownlinkFormatExtendedSquitterNonTransponder ||
		downlinkFormat == DownlinkFormatMilitaryExtendedSquitter) &&
		len(m) == extendedSquitterBitLength/8 {
		return ExtendedSquitter{ModeS: m}, nil
	}

	if (downlinkFormat == DownlinkFormatLongAirAirSurveillance ||
		downlinkFormat == DownlinkFormatCommBWithAltitudeReply ||
		downlinkFormat == DownlinkFormatCommBWithIdentityReply ||
		downlinkFormat == DownlinkFormatCommDExtendedLengthMessage) &&
		len(m) == extendedSquitterBitLength/8 {
		return LongMessage{ModeS: m}, nil
	}

	if (downlinkFormat == DownlinkFormatShortAirAirSurveillance ||
		downlinkFormat == DownlinkFormatAltitudeReply ||
		downlinkFormat == DownlinkFormatIdentityReply ||
		downlinkFormat == DownlinkFormatAllCallReply) &&
		len(m) == shortSquitterBitLength/8 {
		return ShortMessage{ModeS: m}, nil
	}

	return nil, fmt.Errorf("DF:%d / len:%d / msg:%s / err:%w", downlinkFormat, len(m), m, ErrUnsupportedFormat)
}

// IcaoAddrChecksum is when the message type has the checksum xored with the ICAO address,
// this function extracts the ICAO address.
func (m ModeS) IcaoAddrChecksum() ICAOAddr {
	downlinkFormat := m.DownlinkFormat()

	if downlinkFormat == DownlinkFormatShortAirAirSurveillance || // Short air-air surveillance (0)
		downlinkFormat == DownlinkFormatAltitudeReply || // Surveillance, altitude reply (4)
		downlinkFormat == DownlinkFormatIdentityReply || // Surveillance, identity reply (5)
		downlinkFormat == DownlinkFormatLongAirAirSurveillance || // Long air-air survillance (16)
		downlinkFormat == DownlinkFormatCommBWithAltitudeReply || // Comm-B, altitude request (20)
		downlinkFormat == DownlinkFormatCommBWithIdentityReply || // Comm-B, identity request (21)
		downlinkFormat == DownlinkFormatCommDExtendedLengthMessage { // Comm-D ELM (24)
		remainder := binary.ChecksumSquitter(m[:len(m)-3])

		return ICAOAddr(uint32(m[len(m)-1]^byte(remainder)) |
			(uint32(m[len(m)-2]^byte(remainder>>8)) << 8) | //nolint: gomnd
			(uint32(m[len(m)-3]^byte(remainder>>16)) << 16)) //nolint: gomnd
	}

	return 0
}

// CheckSum checks the integrity of the message.
func (m ModeS) CheckSum() error {
	downlinkFormat := m.DownlinkFormat()
	if downlinkFormat == DownlinkFormatAllCallReply || // All call reply (11)
		downlinkFormat == DownlinkFormatExtendedSquitter || // Extended squitter (17)
		downlinkFormat == DownlinkFormatExtendedSquitterNonTransponder { // Extended squitter non transponder (18)
		remainder := binary.ChecksumSquitter(m[:len(m)-3])

		if remainder != m.ParityInterrogator() {
			return fmt.Errorf("%w: found %03X, computed %03X", ErrWrongCRC, m.ParityInterrogator(), remainder)
		}
	}

	return nil
}

// String implements the Stringer interface.
func (m ModeS) String() string {
	return strings.ToUpper(hex.EncodeToString(m))
}
