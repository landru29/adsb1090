package model

import "github.com/landru29/adsb1090/internal/binary"

//       ┏━━━━┓
//       ┃ 4  ┃  Altitude reply
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━━━━━┓
//       ┃ DF | FS | DR | UM | AC | Parity ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┈┈┈┨
//       ┃ 5  |  3 |  5 |  6 | 13 |   24   ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━━━━━┛
//       0    5    8    13   19   32
//
//       ┏━━━━┓
//       ┃ AC ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┓
//       ┃ C1 | A1 | C2 | A2 | C4 | A4 | M  | B1 | Q  | B2 | D2 | B4 | D4 ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┨
//       ┃ 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┛
//
//
//
//       ┏━━━━┓
//       ┃ 5  ┃  Identity reply
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━━━━━┓
//       ┃ DF | FS | DR | UM | ID | Parity ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┈┈┈┨
//       ┃ 5  |  3 |  5 |  6 | 13 |   24   ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━━━━━┛
//       0    5    8    13   19   32
//
//       ┏━━━━┓
//       ┃ ID ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┓
//       ┃ C1 | A1 | C2 | A2 | C4 | A4 | XX | B1 | D1 | B2 | D2 | B4 | D4 ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┨
//       ┃ 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┛

// SurveillanceReplyWithAltitude is the aircraft identification (4).
type SurveillanceReplyWithAltitude struct {
	ShortMessage
}

// FlightStatus is the flight status.
func (r SurveillanceReplyWithAltitude) FlightStatus() FlightStatus {
	return FlightStatus(r.ShortMessage.ModeS[0] & 0x07) //nolint: gomnd
}

// Altitude is the aircraft altitude.
func (r SurveillanceReplyWithAltitude) Altitude() float64 { //nolint: cyclop
	return altitudeFrom13Bits(uint16(binary.ReadBits(r.ShortMessage.ModeS, 19, 13))) //nolint: gomnd
}

// SurveillanceReplyWithIdentification is the aircraft identification (21).
type SurveillanceReplyWithIdentification struct {
	ShortMessage
}

// FlightStatus is the flight status.
func (c SurveillanceReplyWithIdentification) FlightStatus() FlightStatus {
	return FlightStatus(c.ShortMessage.ModeS[0] & 0x07) //nolint: gomnd
}

// Identity is the aircraft identity.
func (c SurveillanceReplyWithIdentification) Identity() Squawk {
	return IdentityFrom12Bits(uint16(binary.ReadBits(c.ShortMessage.ModeS, 19, 13))) //nolint: gomnd
}
