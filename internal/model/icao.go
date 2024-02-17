package model

import (
	"strconv"
	"strings"

	"github.com/landru29/adsb1090/internal/errors"
)

const (
	// ErrWrongICAO is when trying to unmarshal the wrong value.
	ErrWrongICAO errors.Error = "wrong ICAO address"

	minICAOsizeJSON = 2
)

// ICAOAddr is the ICAO aircraft address.
type ICAOAddr uint32

// ParseICAOAddr parses from hexadecimal.
func ParseICAOAddr(str string) (ICAOAddr, error) {
	value, err := strconv.ParseUint(str, 16, 32)

	return ICAOAddr(value), err
}

// MarshalJSON implements the json.Marshaler interface.
func (a ICAOAddr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *ICAOAddr) UnmarshalJSON(data []byte) error {
	if len(data) < minICAOsizeJSON || data[0] != '"' || data[len(data)-1] != '"' {
		return ErrWrongICAO
	}

	value, err := ParseICAOAddr(string(data[1 : len(data)-1]))

	*a = value

	return err
}

// String implements the Stringer interface.
func (a ICAOAddr) String() string {
	return strings.ToUpper(strconv.FormatUint(uint64(a), 16))
}

// Set implements the pflag.Value interface.
func (a *ICAOAddr) Set(str string) error {
	val, err := ParseICAOAddr(str)

	*a = val

	return err
}

// Type implements the pflag.Value interface.
func (a ICAOAddr) Type() string {
	return "OACI addr"
}
