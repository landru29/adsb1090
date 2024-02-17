package model

// ┏━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┓
// ┃ DF  |        Long Message (extended squitter, other)              | Parity ┃
// ┠┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┨
// ┃  5  |                             83                              |   24   ┃
// ┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┛
//

const longMessageName = "long message"

// LongMessage is an extended squitter message.
type LongMessage struct {
	ModeS
}

// AircraftAddress implements the Squitter interface.
func (l LongMessage) AircraftAddress() ICAOAddr {
	return 0
}

// Message is the extended squitter message.
func (l LongMessage) Message() []byte {
	return l.ModeS[4:]
}

// Name implements the Squitter interface.
func (l LongMessage) Name() string {
	return longMessageName
}
