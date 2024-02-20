package model

import (
	"fmt"
	"strings"
	"time"
)

// verticalRate := int64(aircraft.Message.VertRate-1)
//               * 64
//               * map[bool]int64{true: -1, false: 1}[aircraft.Message.VertRateNegative]

// Aircraft is an aircraft description.
type Aircraft struct {
	Identification     string         `json:"ident"`
	CurrentOperation   string         `json:"currentOperation"`
	IcaoAddress        ICAOAddr       `json:"icaoAddress"`
	Altitude           float64        `json:"altitude,omitempty"`
	Position           *Position      `json:"position,omitempty"`
	Flight             string         `json:"flight"` /* Flight number */
	FlightStatus       *FlightStatus  `json:"flightStatus,omitempty"`
	Addr               ICAOAddr       `json:"icao"`                  /* ICAO address */
	GroundSpeed        *float64       `json:"groundSpeed,omitempty"` /* Velocity computed from EW and NS components. */
	AirSpeed           *float64       `json:"airSpeed,omitempty"`    /* Velocity computed from EW and NS components. */
	Track              *float64       `json:"track,omitempty"`       /* Angle of flight. */
	TrueAirSpeed       bool           `json:"trueAirSpeed"`
	BaroVerticalRate   bool           `json:"barometricVerticalRate"`
	DeltaBarometric    int16          `json:"deltaBaro"`
	Identity           Squawk         `json:"identity"`         /* 13 bits identity (from transponder). */
	LastUpdate         time.Time      `json:"lastUpdate"`       /* Time at which the last packet was received. */
	LastFlightStatus   int            `json:"lastFlightStatus"` /* Flight status for DF4,5,20,21 */
	LastDownlinkFormat DownlinkFormat `json:"downlinkFormat"`   /* Downlink format # */
	VerticalRate       int64          `json:"verticalRate"`
	Category           string         `json:"category"`
	Registration       string         `json:"registration"`
	ManufacturerName   string         `json:"manufacturerName"`
	Model              string         `json:"model"`
	Operator           string         `json:"operator"`
	Owner              string         `json:"owner"`
	Built              *time.Time     `json:"built,omitempty"`
	LastType           TypeCode
	LastSubType        SubTypeCode
}

// String implements the Stringer interface.
func (a Aircraft) String() string {
	fields := []string{
		fmt.Sprintf("Ident:     %s", a.Identification),
		fmt.Sprintf("Reg:       %s", a.Registration),
		fmt.Sprintf("Model:     %s", a.Model),
		fmt.Sprintf("Operator:  %s", a.Operator),
		fmt.Sprintf("Addr:      %06X", a.Addr),
		fmt.Sprintf("Category:  %s", a.Category),
		fmt.Sprintf("flight:    %s", a.Flight),
		fmt.Sprintf("altitude:  %f", a.Altitude),
		fmt.Sprintf("seen:      %s", a.LastUpdate.Format(time.RFC3339)),
	}

	if a.GroundSpeed != nil {
		fields = append(fields,
			fmt.Sprintf("speed G:   %f", *a.GroundSpeed),
		)
	}

	if a.AirSpeed != nil {
		fields = append(fields,
			fmt.Sprintf("speed A:   %f", *a.AirSpeed),
		)
	}

	if a.Track != nil {
		fields = append(fields,
			fmt.Sprintf("track:     %f", *a.Track),
		)
	}

	if a.Position != nil {
		fields = append(fields,
			fmt.Sprintf("lat:      %f", a.Position.Latitude),
			fmt.Sprintf("lng:      %f", a.Position.Longitude),
		)
	}

	if a.CurrentOperation != "" {
		fields = append(fields,
			fmt.Sprintf("Operation: %s", a.CurrentOperation),
		)
	}

	return strings.Join(fields, "\n")
}

// Emergency ...
func (a Aircraft) Emergency() bool {
	return (a.LastDownlinkFormat == 4 || a.LastDownlinkFormat == 5 || a.LastDownlinkFormat == 21) &&
		(a.Identity == SquawkHijacker || a.Identity == SquawkRadioFailure || a.Identity == SquawkMayday)
}

// Alert ...
func (a Aircraft) Alert() bool {
	return (a.LastDownlinkFormat == 4 || a.LastDownlinkFormat == 5 || a.LastDownlinkFormat == 21) &&
		(a.LastFlightStatus == 2 || a.LastFlightStatus == 3 || a.LastFlightStatus == 4)
}

// Ground ...
func (a Aircraft) Ground() bool {
	return (a.LastDownlinkFormat == 4 || a.LastDownlinkFormat == 5 || a.LastDownlinkFormat == 21) &&
		(a.LastFlightStatus == 1 || a.LastFlightStatus == 3)
}

// Indent ...
func (a Aircraft) Indent() bool {
	return (a.LastDownlinkFormat == 4 || a.LastDownlinkFormat == 5 || a.LastDownlinkFormat == 21) &&
		(a.LastFlightStatus == 4 || a.LastFlightStatus == 5)
}
