package binary

import "sync"

var (
	crcTables map[uint32]*CRCTable //nolint: gochecknoglobals

	tableAccess sync.Mutex //nolint: gochecknoglobals
)

// CRCTable is a set of XOR masks to compute checksums.
type CRCTable [256]uint32

const (
	crcSize = 256
)

// PolynomCRC is 25 bits (1111111111111010000001001).
// x^24 + x^23 +x^22 +x^21 +x^20 +x^19 +x^18 +x^17 +x^16 +x^15 +x^14 +x^13 +x^12 +x^10 +x^3 + 1
// corresponding to binary: 1111111111111010000001001.
const PolynomCRC uint32 = 0x01fff409

// ChecksumSquitter computes the checksum.
func ChecksumSquitter(input []byte) uint32 {
	table, found := crcTables[PolynomCRC]
	if !found {
		table = MakeCRCTable(PolynomCRC)
	}

	return table.Checksum24(input, len(input)*8) //nolint: gomnd
}

func checksum24(input []byte, polynom uint32) uint32 {
	data := make([]byte, len(input)+3) //nolint: gomnd
	copy(data, input)

	for idx := 0; idx < len(input)*8; idx++ {
		bitIdx := byte(idx % 8)      //nolint: gomnd
		byteIdx := idx / 8           //nolint: gomnd
		mask := byte(0x80) >> bitIdx //nolint: gomnd

		if data[byteIdx]&mask != 0 {
			val := ReadBits(data, uint64(idx), 25)                //nolint: gomnd
			WriteBits(data, val^uint64(polynom), uint64(idx), 25) //nolint: gomnd
		}
	}

	return uint32(ReadBits(data, uint64(len(input)*8), 24)) //nolint: gomnd
}

// Checksum24 makes the checksum.
func (t CRCTable) Checksum24(input []uint8, bitLength int) uint32 {
	var crc uint32

	offset := crcSize - bitLength

	for idx := 0; idx < bitLength; idx++ {
		byteIdx := idx / 8                  //nolint: gomnd
		bitIdx := idx % 8                   //nolint: gomnd
		bitmask := uint8(1) << (7 - bitIdx) //nolint: gomnd

		/* If bit is set, xor with corresponding table entry. */
		if (input[byteIdx] & bitmask) != 0 {
			crc ^= t[idx+offset]
		}
	}

	return crc
}

// MakeCRCTable generate a table of XOR mask to compute checksums.
func MakeCRCTable(poly uint32) *CRCTable {
	tableAccess.Lock()
	defer tableAccess.Unlock()

	if crcTables == nil {
		crcTables = map[uint32]*CRCTable{}
	}

	table := CRCTable{}

	for idx := 0; idx < crcSize; idx++ {
		data := make([]byte, 32) //nolint: gomnd

		byteIdx := idx / 8 //nolint: gomnd
		bitIdx := idx % 8  //nolint: gomnd

		data[byteIdx] = 1 << (7 - bitIdx) //nolint: gomnd

		table[idx] = checksum24(data, poly)
	}

	crcTables[poly] = &table

	return &table
}
