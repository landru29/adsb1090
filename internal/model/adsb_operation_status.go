package model

//       ┏━━━━━┓
//       ┃ 1-4 ┃
//       ┣━━━━━╇━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
//       ┃ TC  | ST |                                                            ┃
//       ┠┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
//       ┃ 5   |  3 |                         48                                 ┃
//       ┗━━━━━╈━━━━╇━━━━━┯━━━━┯━━━━━┯━━━━━━┯━━━━━━┯━━━━━┯━━━━━┯━━━━━┯━━━━━┯━━━━━┫
//             ┃ =0 | CC  | OM | Ver | NICs | NACp | BAQ | SIL | BAI | HRD | Res ┃
//             ┗━━━━╅┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┨
//                  ┃ 16  | 16 |  3  |  1   |  4   |  2  |  2  |  1  |  1  |  2  ┃
//             ┏━━━━╇━━━━━┷━━━━┷━━━━━┷━━━━━━┷━━━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━┫
//             ┃ =1 | CC  | OM | Ver | NICs | NACp | Res | SIL | BAI | HRD | Res ┃
//             ┗━━━━╅┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┨
//                  ┃ 16  | 16 |  3  |  1   |  4   |  2  |  2  |  1  |  1  |  2  ┃
//                  ┗━━━━━┷━━━━┷━━━━━┷━━━━━━┷━━━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━┛

const operationStatusName = "operation status"

// OperationStatus is the operation status.
type OperationStatus struct {
	ExtendedSquitter
}

// Name implements the Message interface.
func (o OperationStatus) Name() string {
	return operationStatusName
}

// String implements the Stringer interface.
// It's the name of the current operation.
func (o OperationStatus) String() string {
	return "" //TODO: Not implemented
}

// Message is the data byte.
func (o OperationStatus) Message() []byte {
	return LongMessage(o.ExtendedSquitter).Message()
}
