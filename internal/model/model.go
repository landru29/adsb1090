// Package model ...
package model

// Namer is a structure with a name.
type Namer interface {
	Name() string
}

// QualifiedMessage is the squitter message.
type QualifiedMessage interface {
	Namer
	AircraftAddress() ICAOAddr
	DownlinkFormat() DownlinkFormat
}

// NamedExtendedSquitter is a received message from extended squitter.
type NamedExtendedSquitter interface {
	Namer
}

// OddEven is type of data that needs odd and even frames.
type OddEven interface {
	OddFrame() bool
}

// UntypeArray removes types from all the elements of an array.
func UntypeArray[T any](data []T) []any {
	output := make([]any, len(data))
	for idx, elt := range data {
		output[idx] = (any)(elt)
	}

	return output
}

// GrayToBinary is the Gray code converter to binary value.
func GrayToBinary(gray uint32) uint32 {
	binary := uint32(0)
	start := false

	for idx := 0; idx < 32; idx++ {
		bit := gray & (1 << (31 - idx)) //nolint: gomnd

		if !start && bit != 0 {
			start = true
			binary = bit

			continue
		}

		if !start {
			continue
		}

		previousBinaryBit := ((binary & (1 << (32 - idx))) >> 1) //nolint: gomnd

		binary |= bit ^ previousBinaryBit
	}

	return binary
}

// func GillhamToBinary(gillham uint32) uint32 {
// 	D4 := (gillham & 0x1 << 10) >> 10
// 	B4 := (gillham & (0x1 << 9)) >> 9
// 	D2 := (gillham & (0x1 << 8)) >> 8
// 	B2 := (gillham & (0x1 << 7)) >> 7
// 	B1 := (gillham & (0x1 << 6)) >> 6
// 	A4 := (gillham & (0x1 << 5)) >> 5
// 	C4 := (gillham & (0x1 << 4)) >> 4
// 	A2 := (gillham & (0x1 << 3)) >> 3
// 	C2 := (gillham & (0x1 << 2)) >> 2
// 	A1 := (gillham & (0x1 << 1)) >> 1
// 	C1 := (gillham & (0x1))

// 	gray := (D2 << 10) +
// 		(D4 << 9) +
// 		(A1 << 8) +
// 		(A2 << 7) +
// 		(A4 << 6) +
// 		(B1 << 5) +
// 		(B2 << 4) +
// 		(B4 << 3) +
// 		(C1 << 2) +
// 		(C2 << 1) +
// 		C4

// 	return GrayToBinary(gray)
// }
