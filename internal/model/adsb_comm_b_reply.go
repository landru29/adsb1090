package model

import "github.com/landru29/adsb1090/internal/binary"

//       ┏━━━━┓
//       ┃ 4  ┃  Altitude reply
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
//       ┃ 5  ┃  Identity reply
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

// CommBReplyWithAltitude is the aircraft identification (4).
type CommBReplyWithAltitude struct {
	LongMessage
}

// FlightStatus is the flight status.
func (c CommBReplyWithAltitude) FlightStatus() FlightStatus {
	return FlightStatus(c.LongMessage.ModeS[0] & 0x07) //nolint: gomnd
}

// Altitude is the aircraft altitude.
func (c CommBReplyWithAltitude) Altitude() float64 { //nolint: cyclop
	idData := binary.ReadBits(c.LongMessage.ModeS, 19, 13) //nolint: gomnd
	bitC1 := (idData & 0x1000) >> 12                       //nolint: gomnd
	bitA1 := (idData & 0x0800) >> 11                       //nolint: gomnd
	bitC2 := (idData & 0x0400) >> 10                       //nolint: gomnd
	bitA2 := (idData & 0x0200) >> 9                        //nolint: gomnd
	bitC4 := (idData & 0x0100) >> 8                        //nolint: gomnd
	bitA4 := (idData & 0x0080) >> 7                        //nolint: gomnd
	bitM := (idData & 0x0040) >> 6                         //nolint: gomnd
	bitB1 := (idData & 0x0020) >> 5                        //nolint: gomnd
	bitQ := (idData & 0x0010) >> 4                         //nolint: gomnd
	bitB2 := (idData & 0x0008) >> 3                        //nolint: gomnd
	bitD2 := (idData & 0x0004) >> 2                        //nolint: gomnd
	bitB4 := (idData & 0x0002) >> 1                        //nolint: gomnd
	bitD4 := idData & 0x0001                               //nolint: gomnd

	if bitC1 == 0 &&
		bitA1 == 0 &&
		bitC2 == 0 &&
		bitA2 == 0 &&
		bitC4 == 0 &&
		bitA4 == 0 &&
		bitM == 0 &&
		bitB1 == 0 &&
		bitQ == 0 &&
		bitB2 == 0 &&
		bitD2 == 0 &&
		bitB4 == 0 &&
		bitD4 == 0 {
		return -1
	}

	if bitM == 1 {
		meters := (idData & 0x27) + ((idData & 0xf80) >> 1) //nolint: gomnd

		return float64(meters) * meterToFeet
	}

	if bitQ == 1 {
		feets := (idData & 0x0f) + (bitB1 << 4) + ((idData & 0xf80) >> 2) //nolint: gomnd

		return float64(feets)*25 - 1000 //nolint: gomnd
	}

	// bitM == 0, bitQ == 0 not implemented.
	return -1
}

// CommBReplyWithIdentification is the aircraft identification (5).
type CommBReplyWithIdentification struct {
	LongMessage
}

// FlightStatus is the flight status.
func (c CommBReplyWithIdentification) FlightStatus() FlightStatus {
	return FlightStatus(c.LongMessage.ModeS[0] & 0x07) //nolint: gomnd
}

// Identity is the aircraft identity.
func (c CommBReplyWithIdentification) Identity() uint16 {
	idData := binary.ReadBits(c.LongMessage.ModeS, 19, 13) //nolint: gomnd
	bitC1 := uint16((idData & 0x1000) >> 12)               //nolint: gomnd
	bitA1 := uint16((idData & 0x0800) >> 11)               //nolint: gomnd
	bitC2 := uint16((idData & 0x0400) >> 10)               //nolint: gomnd
	bitA2 := uint16((idData & 0x0200) >> 9)                //nolint: gomnd
	bitC4 := uint16((idData & 0x0100) >> 8)                //nolint: gomnd
	bitA4 := uint16((idData & 0x0080) >> 7)                //nolint: gomnd
	bitB1 := uint16((idData & 0x0020) >> 5)                //nolint: gomnd
	bitD1 := uint16((idData & 0x0010) >> 4)                //nolint: gomnd
	bitB2 := uint16((idData & 0x0008) >> 3)                //nolint: gomnd
	bitD2 := uint16((idData & 0x0004) >> 2)                //nolint: gomnd
	bitB4 := uint16((idData & 0x0002) >> 1)                //nolint: gomnd
	bitD4 := uint16(idData & 0x0001)                       //nolint: gomnd

	return (bitD1 + (bitD2 << 1) + (bitD4 << 2)) + //nolint: gomnd
		10*(bitC1+(bitC2<<1)+(bitC4<<2)) + //nolint: gomnd
		100*(bitB1+(bitB2<<1)+(bitB4<<2)) + //nolint: gomnd
		1000*(bitA1+(bitA2<<1)+(bitA4<<2)) //nolint: gomnd
}

// FlightStatus is the status of the flight.
type FlightStatus uint8

const (
	// FlightStatusAirborneNoAlertNoSPI : no alert, no SPI, aircraft is airborne.
	FlightStatusAirborneNoAlertNoSPI FlightStatus = 0
	// FlightStatusGroundNoAlertNoSPI : no alert, no SPI, aircraft is on-ground.
	FlightStatusGroundNoAlertNoSPI FlightStatus = 1
	// FlightStatusAirborneAlertNoSPI : alert, no SPI, aircraft is airborne.
	FlightStatusAirborneAlertNoSPI FlightStatus = 2
	// FlightStatusGroundAlertNoSPI : alert, no SPI, aircraft is on-ground.
	FlightStatusGroundAlertNoSPI FlightStatus = 3
	// FlightStatusAlertSPI : alert, SPI, aircraft is airborne or on-ground.
	FlightStatusAlertSPI FlightStatus = 4
	// FlightStatusNoAlertSPI : no alert, SPI, aircraft is airborne or on-ground.
	FlightStatusNoAlertSPI FlightStatus = 5
)

// String implements the Stringer interface.
func (f FlightStatus) String() string {
	switch f {
	case FlightStatusAirborneNoAlertNoSPI:
		return "no alert, no SPI, aircraft is airborne"
	case FlightStatusGroundNoAlertNoSPI:
		return "no alert, no SPI, aircraft is on-ground"
	case FlightStatusAirborneAlertNoSPI:
		return "alert, no SPI, aircraft is airborne"
	case FlightStatusGroundAlertNoSPI:
		return "alert, no SPI, aircraft is on-ground"
	case FlightStatusAlertSPI:
		return "alert, SPI, aircraft is airborne or on-ground"
	case FlightStatusNoAlertSPI:
		return "no alert, SPI, aircraft is airborne or on-ground"
	default:
		return "invalid"
	}
}
