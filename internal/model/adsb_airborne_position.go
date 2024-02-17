package model

//       ┏━━━━━━━┓
//       ┃ 8-18  ┃
//       ┃ 20-23 ┃
//       ┣━━━━━━━╇━━━━┯━━━━━┯━━━━━┯━━━┯━━━┯━━━━━━━━━┯━━━━━━━━━┓
//       ┃  TC   | SS | SAF | ALT | T | F | LAT-CPR | LON-CPR ┃
//       ┠┈┈┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┼┈┈┈┼┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┨
//       ┃   5   |  2 |  1  |  12 | 1 | 1 |    17   |   17    ┃
//       ┗━━━━━━━┷━━━━┷━━━━━┷━━━━━┷━━━┷━━━┷━━━━━━━━━┷━━━━━━━━━┛

const (
	airbornePositionName = "airborne position"
	meterToFeet          = 3.28084
)

// AirbornePosition is the surface position.
type AirbornePosition struct {
	ExtendedSquitter
}

// DecodePosition decodes the current position with another frame.
func (p AirbornePosition) DecodePosition(other AirbornePosition) (*Position, error) {
	return DecodePosition(p, other)
}

// Altitude is the aircraft altitude.
func (p AirbornePosition) Altitude() float64 {
	message := p.Message()

	typeCode := p.TypeCode()

	encodedAltitude := (uint16(message[1]&0x7f) << 5) + uint16((message[2]&0xf8)>>3) //nolint: gomnd

	// barometric Altitude
	if typeCode > 8 && typeCode < 19 {
		Qbit := encodedAltitude & 0x10 //nolint: gomnd

		encodedAltitude = encodedAltitude&0x0f + ((encodedAltitude & 0xff0) >> 1) //nolint: gomnd

		if Qbit == 1 {
			return float64(encodedAltitude)*25 - 1000 //nolint: gomnd
		}

		// In the case where the altitude is higher than 50175 feet,
		// a 100 feet increment is used. In this situation, the Q bit is set to 0,
		// and the rest of the bits are encoded using Gray code [Doran 2007].
		// Not yet implemented due to lack of samples (for testing purpose).
		return -1
	}

	// GNSS altitude (in meters)
	if typeCode > 19 && typeCode < 23 {
		return float64(encodedAltitude) * meterToFeet
	}

	return 0
}

// Name implements the Message interface.
func (p AirbornePosition) Name() string {
	return airbornePositionName
}

// SurveillanceStatus is the surveillance status.
func (p AirbornePosition) SurveillanceStatus() SurveillanceStatus {
	return SurveillanceStatus((p.Message()[0] & 0x6) >> 1) //nolint: gomnd
}

// SingleAntennaFlag defines if the antenna is single or dual.
func (p AirbornePosition) SingleAntennaFlag() bool {
	return map[byte]bool{1: true, 0: false}[p.Message()[0]&0x1]
}

// EncodedAltitude is the encoded altitude.
func (p AirbornePosition) EncodedAltitude() uint16 {
	message := p.Message()

	return (uint16(message[1]) << 4) | (uint16(message[2]) >> 4) //nolint: gomnd
}

// TimeUTC define whether the time is UTC or not.
func (p AirbornePosition) TimeUTC() bool {
	return map[byte]bool{1: true, 0: false}[(p.Message()[2]>>3)&0x1] //nolint: gomnd
}

// OddFrame defines if the frame is odd or even.
func (p AirbornePosition) OddFrame() bool {
	return map[byte]bool{1: true, 0: false}[(p.Message()[2]>>2)&0x1] //nolint: gomnd
}

// EncodedLatitude is the encoded latitude.
func (p AirbornePosition) EncodedLatitude() uint32 {
	message := p.Message()

	return ((uint32(message[2]) & 0x3) << 15) | //nolint: gomnd
		(uint32(message[3]) << 7) | //nolint: gomnd
		(uint32(message[4]) >> 1)
}

// EncodedLongitude is the encoded longitude.
func (p AirbornePosition) EncodedLongitude() uint32 {
	message := p.Message()

	return ((uint32(message[4]) & 0x1) << 16) | //nolint: gomnd
		(uint32(message[5]) << 8) | //nolint: gomnd
		uint32(message[6])
}

// Baro defines whether the message is based on a baro altitude or GNSS altitude.
func (p AirbornePosition) Baro() bool {
	return p.TypeCode() == TypeCodeAirbornePositionBaroAltitude
}

// Message is the data byte.
func (p AirbornePosition) Message() []byte {
	return LongMessage(p.ExtendedSquitter).Message()
}
