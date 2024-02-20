package config

import (
	"fmt"

	"github.com/landru29/adsb1090/internal/serialize/nmea"
)

// Vessel is a vessel type for NMEA purpose.
type Vessel nmea.VesselType

// String implements the pflag.Value interface.
func (v Vessel) String() string {
	return map[nmea.VesselType]string{
		nmea.VesselTypeAircraft:   "aircraft",
		nmea.VesselTypeHelicopter: "helicopter",
	}[nmea.VesselType(v)]
}

// Set implements the pflag.Value interface.
func (v *Vessel) Set(str string) error {
	vesselType, ok := map[string]nmea.VesselType{
		"aircraft":   nmea.VesselTypeAircraft,
		"helicopter": nmea.VesselTypeHelicopter,
	}[str]
	if !ok {
		return fmt.Errorf("unknow vessel type %s", str)
	}

	*v = Vessel(vesselType)

	return nil
}

// Type implements the pflag.Value interface.
func (v Vessel) Type() string {
	return "vessel type"
}
