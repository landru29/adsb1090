package nmea

import (
	"bytes"

	"github.com/landru29/adsb1090/internal/model"
)

const (
	speedOverGroundScale = 10
)

// VesselType is a type of vessel.
type VesselType int

const (
	// VesselTypeAircraft is an aircraft.
	VesselTypeAircraft = iota
	// VesselTypeHelicopter is a helicopter.
	VesselTypeHelicopter
)

// Serializer is the nmea serializer.
type Serializer struct {
	mmsiVessel VesselType
	mid        uint16
}

// New is a new NMEA serializer.
func New(mmsiVessel VesselType, mid uint16) *Serializer {
	return &Serializer{
		mmsiVessel: mmsiVessel,
		mid:        mid,
	}
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(planes ...any) ([]byte, error) {
	output := [][]byte{}

	for _, ac := range planes {
		switch aircraft := ac.(type) {
		case model.Aircraft:
			data, err := s.Serialize(&aircraft)
			if err != nil {
				return nil, err
			}

			output = append(output, data)

		case *model.Aircraft:
			if aircraft != nil && aircraft.Position != nil {
				fields, err := s.fieldFromAircraft(aircraft)
				if err != nil {
					return nil, err
				}

				output = append(output, []byte(fields.String()), []byte("\n"))
			}
		case []model.Aircraft:
			data, err := s.Serialize(model.UntypeArray(aircraft)...)
			if err != nil {
				return nil, err
			}

			output = append(output, data)
		case []*model.Aircraft:
			data, err := s.Serialize(model.UntypeArray(aircraft)...)
			if err != nil {
				return nil, err
			}

			output = append(output, data)
		}
	}

	return bytes.Join(output, []byte("\n")), nil
}

func (s Serializer) fieldFromAircraft(aircraft *model.Aircraft) (fields, error) {
	currentPayload := payload{
		MMSI:             s.MMSI(aircraft.Addr),
		Longitude:        aircraft.Position.Longitude,
		Latitude:         aircraft.Position.Latitude,
		PositionAccuracy: true,
		NavigationStatus: navigationStatusAground,
	}

	if aircraft.AirSpeed != nil {
		currentPayload.SpeedOverGround = *aircraft.AirSpeed / speedOverGroundScale
	}

	if aircraft.GroundSpeed != nil {
		currentPayload.SpeedOverGround = *aircraft.GroundSpeed / speedOverGroundScale
	}

	if aircraft.Track != nil {
		currentPayload.CourseOverGround = *aircraft.Track
		currentPayload.TrueHeading = uint16(*aircraft.Track)
	}

	return currentPayload.Fields()
}

// MimeType implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "application/nmea"
}

// MMSI ...
func (s Serializer) MMSI(addr model.ICAOAddr) uint32 {
	out := uint32(s.mid%1000)*10000 + 10000000 //nolint: gomnd

	switch s.mmsiVessel {
	case VesselTypeAircraft:
		out += 1000
	case VesselTypeHelicopter:
		out += 5000
	}

	return out + uint32(addr)%1000
}

// String implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "nmea"
}
