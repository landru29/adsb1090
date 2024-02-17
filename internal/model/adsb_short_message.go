package model

// ┏━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┓
// ┃ DF  |                        Short Message                        | Parity ┃
// ┠┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┨
// ┃  5  |                             27                              |   24   ┃
// ┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┛
//

const shortMessageName = "short message"

// ShortMessage is a short squitter message.
type ShortMessage struct {
	ModeS
}

// AircraftAddress implements the Squitter interface.
func (s ShortMessage) AircraftAddress() ICAOAddr {
	return s.IcaoAddrChecksum()
}

// Name implements the Squitter interface.
func (s ShortMessage) Name() string {
	return shortMessageName
}
