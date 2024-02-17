package model

import "math"

//       ┏━━━━┓
//       ┃ 19 ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━━┯━━━━━┯━━━━━━━━━━━━━━━━━┯━━━━━━━┯━━━━━┯━━━━┯━━━━━┯━━━━━━┯━━━━━━┓
//       ┃ TC | ST | IC | IFR | NUC | Specific Fields | VrSrc | Svr | VR | Res | SDif | dAlt ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┨
//       ┃ 5  |  3 |  1 |  1  |  3  |       22        |   1   |  1  | 9  |  2  |  1   |  7   ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━━┷━━━━━┷━━━━━━━━━━━━━━━━━┷━━━━━━━┷━━━━━┷━━━━┷━━━━━┷━━━━━━┷━━━━━━┛
//       0    5    8    9     10    13                35      36    37   46    48     49     56
//
//
//  Specific fields:
//  ---------------
//
//  ST=1,2
//       ┏━━━━━┯━━━━━┯━━━━━┯━━━━━┓
//       ┃ Dew | Vew | Dns | Vns ┃
//       ┠┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┨
//       ┃ 1   |  10 |  1  | 10  ┃
//       ┗━━━━━┷━━━━━┷━━━━━┷━━━━━┛
//       13    14    24    25
//
//  ST=3,4
//       ┏━━━━┯━━━━━┯━━━┯━━━━┓
//       ┃ SH | HDG | T | AS ┃
//       ┠┈┈┈┈┼┈┈┈┈┈┼┈┈┈┼┈┈┈┈┨
//       ┃ 1  |  10 | 1 | 10 ┃
//       ┗━━━━┷━━━━━┷━━━┷━━━━┛
//       13   14    24  25
//
//  0    0 - 7
//  1    8 - 15
//  2   16 - 23
//  3   24 - 31
//  4   32 - 39
//  5   40 - 47
//  6   48 - 55

const airborneVelocityName = "airborne velocity"

// AirborneVelocity is the surface position.
type AirborneVelocity struct {
	ExtendedSquitter
}

func (v AirborneVelocity) subType() airborneVelocitySubType {
	return airborneVelocitySubType(v.Message()[0] & 0x07) //nolint: gomnd
}

// IsBaroVerticalRate indicates the source of the altitude measurements (GNSS altitude or barometric altitude).
func (v AirborneVelocity) IsBaroVerticalRate() bool {
	return v.Message()[4]&0x10 != 0
}

// Name implements the Message interface.
func (v AirborneVelocity) Name() string {
	return airborneVelocityName
}

func (v AirborneVelocity) subTypeFields() (bool, bool, int16, int16) {
	one := (v.Message()[1] & 0x04) >> 2                                             //nolint: gomnd
	two := (uint16(v.Message()[1]&0x03) << 8) + uint16(v.Message()[2])              //nolint: gomnd
	three := (v.Message()[3] & 0x80) >> 7                                           //nolint: gomnd
	four := (uint16(v.Message()[3]&0x7f) << 3) + (uint16(v.Message()[4]&0xe0) >> 5) //nolint: gomnd

	return one == 1, three == 1, int16(two), int16(four)
}

// IsTrueAirSpeed checks if the airspeed is estimated.
func (v AirborneVelocity) IsTrueAirSpeed() bool {
	subtype := v.subType()

	if subtype == 3 || subtype == 4 {
		return (v.Message()[3]&0x80)>>7 != 0 //nolint: gomnd
	}

	return false
}

// VerticalRate is the vertival speed.
func (v AirborneVelocity) VerticalRate() int64 {
	absoluteRate := ((uint32(v.Message()[4]&0x07) << 6) + (uint32(v.Message()[5]&0xfc) >> 2) - 1) * 64 //nolint: gomnd

	if v.Message()[4]&0x08 != 0 {
		return -int64(absoluteRate)
	}

	return int64(absoluteRate)
}

// DeltaBarometric is the GNSS and barometric altitudes difference.
func (v AirborneVelocity) DeltaBarometric() int16 {
	data := v.Message()[6] & 0x7f //nolint: gomnd
	if data == 0 {
		return 0
	}

	absoluteDelta := 25 * (int16(data) - 1) //nolint: gomnd

	if v.Message()[6]&0x80 != 0 {
		return -absoluteDelta
	}

	return absoluteDelta
}

func (v AirborneVelocity) groundSpeedVector(factor int16) (float64, float64) {
	dew, dns, vew, vns := v.subTypeFields()

	speedX := vew - 1
	speedY := vns - 1

	if dew {
		speedX = -factor * speedX
	}

	if dns {
		speedY = -factor * speedY
	}

	return float64(speedX), float64(speedY)
}

func (v AirborneVelocity) airSpeedPolarVector(factor int16) (float64, float64) {
	headingStatusBit, _, hdg, airSpeed := v.subTypeFields()

	heading := float64(-1)

	if headingStatusBit {
		heading = float64(360.0/1024.0) * float64(hdg) //nolint: gomnd
	}

	return float64(factor * (airSpeed - 1)), heading
}

// IsGroundSpeed gives the type of the speed.
func (v AirborneVelocity) IsGroundSpeed() bool {
	subType := v.subType()
	if subType == 1 || subType == 2 {
		return true
	}

	return false
}

// Speed is the speed value and the heading.
func (v AirborneVelocity) Speed() (float64, float64) {
	message := v.Message()
	if (message[1]&0x07) == 0 && message[2] == 0 && message[3] == 0 && (message[4]&0xe0) == 0 { //nolint: gomnd
		return -1, -1
	}

	subType := v.subType()
	switch subType {
	case 1:
		vx, vy := v.groundSpeedVector(1)

		return math.Sqrt(vx*vx + vy*vy), math.Mod(math.Atan2(vx, vy)*180.0/math.Pi+360, 360) //nolint: gomnd

	case 2: //nolint: gomnd
		vx, vy := v.groundSpeedVector(2) //nolint: gomnd

		return math.Sqrt(vx*vx + vy*vy), math.Mod(math.Atan2(vx, vy)*180.0/math.Pi+360, 360) //nolint: gomnd

	case 3: //nolint: gomnd
		return v.airSpeedPolarVector(1) //nolint: gomnd
	case 4: //nolint: gomnd
		return v.airSpeedPolarVector(4) //nolint: gomnd
	}

	return -1, -1
}

// Message is the data byte.
func (v AirborneVelocity) Message() []byte {
	return LongMessage(v.ExtendedSquitter).Message()
}

type airborneVelocitySubType uint8
