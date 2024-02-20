package model

// FlightStatus is the status of the flight.
type FlightStatus uint8

const (
	// FlightStatusAirborneNoAlertNoSPI : no alert, no SPI, aircraft is airborne.
	FlightStatusAirborneNoAlertNoSPI FlightStatus = 0
	// FlightStatusGroundNoAlertNoSPI : no alert, no SPI, aircraft is on-ground.
	FlightStatusGroundNoAlertNoSPI FlightStatus = 1
	// FlightStatusAirborneAlertNoSPI : alert, no SPI, aircraft is airborne.
	FlightStatusAirborneAlertNoSPI FlightStatus = 2
	// FlightStatusGroundAlertNoSPI : alert, no SPI, aircraft is on-ground.
	FlightStatusGroundAlertNoSPI FlightStatus = 3
	// FlightStatusAlertSPI : alert, SPI, aircraft is airborne or on-ground.
	FlightStatusAlertSPI FlightStatus = 4
	// FlightStatusNoAlertSPI : no alert, SPI, aircraft is airborne or on-ground.
	FlightStatusNoAlertSPI FlightStatus = 5
)

// String implements the Stringer interface.
func (f FlightStatus) String() string {
	switch f {
	case FlightStatusAirborneNoAlertNoSPI:
		return "no alert, no SPI, aircraft is airborne"
	case FlightStatusGroundNoAlertNoSPI:
		return "no alert, no SPI, aircraft is on-ground"
	case FlightStatusAirborneAlertNoSPI:
		return "alert, no SPI, aircraft is airborne"
	case FlightStatusGroundAlertNoSPI:
		return "alert, no SPI, aircraft is on-ground"
	case FlightStatusAlertSPI:
		return "alert, SPI, aircraft is airborne or on-ground"
	case FlightStatusNoAlertSPI:
		return "no alert, SPI, aircraft is airborne or on-ground"
	default:
		return "invalid"
	}
}
