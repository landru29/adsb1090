package model

func IdentityFrom12Bits(idData uint16) Squawk {
	bitC1 := uint16((idData & 0x1000) >> 12) //nolint: gomnd
	bitA1 := uint16((idData & 0x0800) >> 11) //nolint: gomnd
	bitC2 := uint16((idData & 0x0400) >> 10) //nolint: gomnd
	bitA2 := uint16((idData & 0x0200) >> 9)  //nolint: gomnd
	bitC4 := uint16((idData & 0x0100) >> 8)  //nolint: gomnd
	bitA4 := uint16((idData & 0x0080) >> 7)  //nolint: gomnd
	bitB1 := uint16((idData & 0x0020) >> 5)  //nolint: gomnd
	bitD1 := uint16((idData & 0x0010) >> 4)  //nolint: gomnd
	bitB2 := uint16((idData & 0x0008) >> 3)  //nolint: gomnd
	bitD2 := uint16((idData & 0x0004) >> 2)  //nolint: gomnd
	bitB4 := uint16((idData & 0x0002) >> 1)  //nolint: gomnd
	bitD4 := uint16(idData & 0x0001)         //nolint: gomnd

	return Squawk((bitD1 + (bitD2 << 1) + (bitD4 << 2)) + //nolint: gomnd
		10*(bitC1+(bitC2<<1)+(bitC4<<2)) + //nolint: gomnd
		100*(bitB1+(bitB2<<1)+(bitB4<<2)) + //nolint: gomnd
		1000*(bitA1+(bitA2<<1)+(bitA4<<2))) //nolint: gomnd
}
