package model

func altitudeFrom13Bits(altitudeData uint16) float64 {
	bitC1 := (altitudeData & 0x1000) >> 12 //nolint: gomnd
	bitA1 := (altitudeData & 0x0800) >> 11 //nolint: gomnd
	bitC2 := (altitudeData & 0x0400) >> 10 //nolint: gomnd
	bitA2 := (altitudeData & 0x0200) >> 9  //nolint: gomnd
	bitC4 := (altitudeData & 0x0100) >> 8  //nolint: gomnd
	bitA4 := (altitudeData & 0x0080) >> 7  //nolint: gomnd
	bitM := (altitudeData & 0x0040) >> 6   //nolint: gomnd
	bitB1 := (altitudeData & 0x0020) >> 5  //nolint: gomnd
	bitQ := (altitudeData & 0x0010) >> 4   //nolint: gomnd
	bitB2 := (altitudeData & 0x0008) >> 3  //nolint: gomnd
	bitD2 := (altitudeData & 0x0004) >> 2  //nolint: gomnd
	bitB4 := (altitudeData & 0x0002) >> 1  //nolint: gomnd
	bitD4 := altitudeData & 0x0001         //nolint: gomnd

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
		meters := (altitudeData & 0x27) + ((altitudeData & 0xf80) >> 1) //nolint: gomnd

		return float64(meters) * meterToFeet
	}

	if bitQ == 1 {
		feets := (altitudeData & 0x0f) + (bitB1 << 4) + ((altitudeData & 0x1f80) >> 2) //nolint: gomnd

		return float64(feets)*25 - 1000 //nolint: gomnd
	}

	// bitM == 0, bitQ == 0 not implemented.
	return -1
}
