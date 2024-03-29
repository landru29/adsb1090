package model

import "github.com/landru29/adsb1090/internal/binary"

//       ┏━━━━┓
//       ┃ 20 ┃  Altitude reply
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━━━━━┓
//       ┃ DF | FS | DR | UM | AC | MB | Parity ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┈┈┈┨
//       ┃ 5  |  3 |  5 |  6 | 13 | 56 |   24   ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━━━━━┛
//       0    5    8    13   19   32   88
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
//       ┃ 21 ┃  Identity reply
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━━━━━┓
//       ┃ DF | FS | DR | UM | ID | MB | Parity ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┈┈┈┨
//       ┃ 5  |  3 |  5 |  6 | 13 | 56 |   24   ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━━━━━┛
//
//       ┏━━━━┓
//       ┃ ID ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┓
//       ┃ C1 | A1 | C2 | A2 | C4 | A4 | XX | B1 | D1 | B2 | D2 | B4 | D4 ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┨
//       ┃ 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  | 1  ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┛

// CommBReplyWithAltitude is the aircraft altitude (20).
type CommBReplyWithAltitude struct {
	LongMessage
}

// FlightStatus is the flight status.
func (c CommBReplyWithAltitude) FlightStatus() FlightStatus {
	return FlightStatus(c.LongMessage.ModeS[0] & 0x07) //nolint: gomnd
}

// Altitude is the aircraft altitude.
func (c CommBReplyWithAltitude) Altitude() float64 { //nolint: cyclop
	return altitudeFrom13Bits(uint16(binary.ReadBits(c.LongMessage.ModeS, 19, 13))) //nolint: gomnd
}

// CommBReplyWithIdentification is the aircraft identification (21).
type CommBReplyWithIdentification struct {
	LongMessage
}

// FlightStatus is the flight status.
func (c CommBReplyWithIdentification) FlightStatus() FlightStatus {
	return FlightStatus(c.LongMessage.ModeS[0] & 0x07) //nolint: gomnd
}

// Identity is the aircraft identity.
func (c CommBReplyWithIdentification) Identity() Squawk {
	return IdentityFrom12Bits(uint16(binary.ReadBits(c.LongMessage.ModeS, 19, 13))) //nolint: gomnd
}
