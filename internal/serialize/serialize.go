// Package serialize describes how to serialize aircraft.
package serialize

import "github.com/landru29/adsb1090/internal/errors"

const (
	// ErrMissingSerializer is when the serializer is missing.
	ErrMissingSerializer errors.Error = "missing serializer"
)

// Serializer is the aircraft serializer.
type Serializer interface {
	Serialize(ac ...any) ([]byte, error)
	MimeType() string
	String() string
}
