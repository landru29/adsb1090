package decoder

import (
	"log/slog"
	"time"

	"github.com/landru29/adsb1090/internal/aircraftdb"
	"github.com/landru29/adsb1090/internal/model"
)

type extendedSquitter struct {
	Identification         []model.Identification
	OperationStatus        []model.OperationStatus
	SurfacePosition        []model.SurfacePosition
	AirbornePosition       []model.AirbornePosition
	AirborneVelocity       []model.AirborneVelocity
	lastPositionIsAirborne bool
}

func buildAircraft(log *slog.Logger, squitters []model.QualifiedMessage, ref aircraftdb.Entry) *model.Aircraft {
	aircraft := model.Aircraft{
		Registration:     ref.Registration,
		ManufacturerName: ref.ManufacturerName,
		Model:            ref.Model,
		Operator:         ref.Operator,
		Owner:            ref.Owner,
		Built:            ref.Built,
		Identity:         0, // TODO get identification
	}

	lastSquitter := squitters[len(squitters)-1]

	aircraft.LastDownlinkFormat = lastSquitter.DownlinkFormat()
	if lastExtendedSquitter, ok := lastSquitter.(model.ExtendedSquitter); ok {
		aircraft.LastType = lastExtendedSquitter.TypeCode()
		aircraft.LastSubType = lastExtendedSquitter.SubTypeCode()
	}

	extendedSquitters := extendedSquitter{}

	for _, genericSquitter := range squitters {
		if extendedSquitter, ok := genericSquitter.(model.ExtendedSquitter); ok {
			decoded, err := extendedSquitter.Decode()
			if err != nil {
				continue
			}

			switch val := decoded.(type) {
			case model.Identification:
				extendedSquitters.Identification = append(extendedSquitters.Identification, val)
			case model.OperationStatus:
				extendedSquitters.OperationStatus = append(extendedSquitters.OperationStatus, val)
			case model.SurfacePosition:
				extendedSquitters.SurfacePosition = append(extendedSquitters.SurfacePosition, val)
				extendedSquitters.lastPositionIsAirborne = false
			case model.AirbornePosition:
				extendedSquitters.AirbornePosition = append(extendedSquitters.AirbornePosition, val)
				extendedSquitters.lastPositionIsAirborne = true
			case model.AirborneVelocity:
				extendedSquitters.AirborneVelocity = append(extendedSquitters.AirborneVelocity, val)
			}
		}

		if shortSquitter, ok := genericSquitter.(model.ShortMessage); ok {
			processShortMessage(log, &aircraft, shortSquitter)
		}
	}

	processExtendedSquitter(log, &aircraft, extendedSquitters)

	return &aircraft
}

func processShortMessage(log *slog.Logger, aircraft *model.Aircraft, message model.ShortMessage) { //nolint: revive,unparam,lll,whitespace,wsl

	// switch squitter.DownlinkFormat() { //nolint: exhaustive
	// case model.DownlinkFormatShortAirAirSurveillance:
	// 	log.Debug("DownlinkFormatShortAirAirSurveillance")

	// case model.DownlinkFormatAltitudeReply:
	// 	log.Debug("DownlinkFormatAltitudeReply")

	// case model.DownlinkFormatIdentityReply:
	// 	log.Debug("DownlinkFormatIdentityReply")

	// case model.DownlinkFormatAllCallReply:
	// 	log.Debug("DownlinkFormatAllCallReply")

	// default:
	// 	return
	// }

	aircraft.LastUpdate = time.Now()
}

func processExtendedSquitter(log *slog.Logger, aircraft *model.Aircraft, squitter extendedSquitter) { //nolint: gocognit,cyclop,lll
	if len(squitter.Identification) > 0 {
		identification := squitter.Identification[len(squitter.Identification)-1]
		aircraft.Category = identification.CategoryString()
		aircraft.Identification = identification.String()
	}

	if !squitter.lastPositionIsAirborne && len(squitter.SurfacePosition) > 0 {
		aircraft.Altitude = squitter.SurfacePosition[len(squitter.SurfacePosition)-1].Altitude()
	}

	if !squitter.lastPositionIsAirborne && len(squitter.SurfacePosition) > 1 {
		lastPosition, otherPosition := frames(squitter.SurfacePosition)

		if otherPosition != nil {
			position, err := lastPosition.DecodePosition(*otherPosition)
			if err != nil {
				log.Debug("surface position error")
			} else {
				aircraft.Position = position
			}
		}
	}

	if squitter.lastPositionIsAirborne && len(squitter.AirbornePosition) > 0 {
		aircraft.Altitude = squitter.AirbornePosition[len(squitter.AirbornePosition)-1].Altitude()
	}

	if squitter.lastPositionIsAirborne && len(squitter.AirbornePosition) > 1 {
		lastPosition, otherPosition := frames(squitter.AirbornePosition)

		if otherPosition != nil {
			position, err := lastPosition.DecodePosition(*otherPosition)
			if err != nil {
				log.Debug("airborne position error")
			} else {
				aircraft.Position = position
			}
		}
	}

	if len(squitter.AirborneVelocity) > 0 { //nolint: nestif
		lastVelocity := squitter.AirborneVelocity[len(squitter.AirborneVelocity)-1]

		speed, heading := lastVelocity.Speed()

		if speed >= 0 {
			if heading > 0 {
				aircraft.Track = &heading
			}

			if lastVelocity.IsGroundSpeed() {
				aircraft.GroundSpeed = &speed
			} else {
				aircraft.AirSpeed = &speed
			}

			aircraft.BaroVerticalRate = lastVelocity.IsBaroVerticalRate()
			aircraft.TrueAirSpeed = lastVelocity.IsTrueAirSpeed()
			aircraft.VerticalRate = lastVelocity.VerticalRate()
			aircraft.DeltaBarometric = lastVelocity.DeltaBarometric()
		}
	}

	if len(squitter.OperationStatus) > 0 {
		aircraft.CurrentOperation = squitter.OperationStatus[len(squitter.OperationStatus)-1].String()
	}

	aircraft.LastUpdate = time.Now()
}

func frames[T model.OddEven](dataSet []T) (T, *T) { //nolint: ireturn
	first := dataSet[len(dataSet)-1]

	return first, findFrame[T](dataSet, !first.OddFrame())
}

func findFrame[T model.OddEven](dataSet []T, odd bool) *T {
	for idx := len(dataSet) - 2; idx >= 0; idx-- { //nolint: gomnd
		if dataSet[idx].OddFrame() == odd {
			return &dataSet[idx]
		}
	}

	return nil
}
