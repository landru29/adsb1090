package model

//       ┏━━━━┓
//       ┃ 31 ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┓
//       ┃ TC | CA | C1 | C2 | C3 | C4 | C5 | C6 | C7 | C8 ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┨
//       ┃ 5  |  3 |  6 |  6 |  6 |  6 |  6 |  6 |  6 |  6 ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┛

const (
	asciiTable         = "#ABCDEFGHIJKLMNOPQRSTUVWXYZ##### ###############0123456789######"
	identificationName = "identification"
)

// Category is the aircraft category.
type Category byte

// Identification is the aircraft identification.
type Identification struct {
	ExtendedSquitter
}

// String implement the Stringer interface.
func (i Identification) String() string {
	message := LongMessage(i.ExtendedSquitter).Message()
	letters := make([]byte, 8) //nolint: gomnd

	letters[0] = asciiTable[message[1]>>2]
	letters[1] = asciiTable[(message[1]&0x3)<<4+message[2]>>4] //nolint: gomnd
	letters[2] = asciiTable[(message[2]&0xf)<<2+message[3]>>6] //nolint: gomnd
	letters[3] = asciiTable[message[3]&0x3f]
	letters[4] = asciiTable[message[4]>>2]
	letters[5] = asciiTable[(message[4]&0x3)<<4+message[5]>>4] //nolint: gomnd
	letters[6] = asciiTable[(message[5]&0xf)<<2+message[6]>>6] //nolint: gomnd
	letters[7] = asciiTable[message[6]&0x3f]

	return string(letters)
}

// Name implements the Message interface.
func (i Identification) Name() string {
	return identificationName
}

// Category is the aircraft category.
func (i Identification) Category() Category {
	return Category(i.Message()[0] & 0x7) //nolint: gomnd
}

// Message is the data byte.
func (i Identification) Message() []byte {
	return LongMessage(i.ExtendedSquitter).Message()
}

// CategoryString is the aircraft category.
func (i Identification) CategoryString() string { //nolint: cyclop
	switch i.Message()[0] {
	case 16 + 1: //nolint: gomnd
		return categorySurfaceEmergencyVehicule
	case 16 + 3: //nolint: gomnd
		return categorySurfaceServiceVehicle
	case 16 + 4, 16 + 5, 16 + 6, 16 + 7: //nolint: gomnd
		return categoryGroundObstruction
	case 24 + 1: //nolint: gomnd
		return categoryLliderSailplane
	case 24 + 2: //nolint: gomnd
		return categoryLighterThanAir
	case 24 + 3: //nolint: gomnd
		return categoryParachutistSkydiver
	case 24 + 4: //nolint: gomnd
		return categoryUltralight
	case 24 + 5: //nolint: gomnd
		return categoryReserved
	case 24 + 6: //nolint: gomnd
		return categoryUnmannedAerialVehicle
	case 24 + 7: //nolint: gomnd
		return categorySpaceOrTransatmosphericVehicle
	case 32 + 1: //nolint: gomnd
		return categoryLighter
	case 32 + 2: //nolint: gomnd
		return categoryMedium1
	case 32 + 3: //nolint: gomnd
		return categoryMedium2
	case 32 + 4: //nolint: gomnd
		return categoryHighVortexAircraft
	case 32 + 5: //nolint: gomnd
		return categoryHeavy
	case 32 + 6: //nolint: gomnd
		return categoryHighPerformance
	case 32 + 7: //nolint: gomnd
		return categoryRotorcraft
	default:
		return categoryNoInformation
	}
}

const (
	categoryNoInformation                  = "No category information"
	categorySurfaceEmergencyVehicule       = "Surface emergency vehicle"
	categorySurfaceServiceVehicle          = "Surface service vehicle"
	categoryGroundObstruction              = "Ground obstruction"
	categoryLliderSailplane                = "Glider, sailplane"
	categoryLighterThanAir                 = "Lighter-than-air"
	categoryParachutistSkydiver            = "Parachutist, skydiver"
	categoryUltralight                     = "Ultralight, hang-glider, paraglider"
	categoryReserved                       = "Reserved"
	categoryUnmannedAerialVehicle          = "Unmanned aerial vehicle"
	categorySpaceOrTransatmosphericVehicle = "Space or transatmospheric vehicle"
	categoryLighter                        = "Light (less than 7000 kg)"
	categoryMedium1                        = "Medium 1 (between 7000 kg and 34000 kg)"
	categoryMedium2                        = "Medium 2 (between 34000 kg to 136000 kg)"
	categoryHighVortexAircraft             = "High vortex aircraft"
	categoryHeavy                          = "Heavy (larger than 136000 kg)"
	categoryHighPerformance                = "High performance (>5 g acceleration) and high speed (>400 kt)"
	categoryRotorcraft                     = "Rotorcraft"
)
