package dns

type OperationCode int // 4 bit

const (
	OperationCodeQuery             = 0 // RFC 1035, 4.1.1
	_                              = 1 // RFC 1035, 4.1.1 (Obsolete)
	OperationCodeStatus            = 2 // RFC 1035, 4.1.1
	_                              = 3 // RFC 1035, 4.1.1 (Unused)
	OperationCodeNotify            = 4 // RFC 1996
	OperationCodeUpdate            = 5 // RFC 2136
	OperationCodeStatefulOperation = 6 // RFC 8490
)

func (opcode OperationCode) String() string {

	switch opcode {
	case OperationCodeQuery:
		return "Query"
	case OperationCodeStatus:
		return "Status"
	case OperationCodeNotify:
		return "Notify"
	case OperationCodeUpdate:
		return "Update"
	case OperationCodeStatefulOperation:
		return "Stateful Operation"
	default:
		return ""
	}

}
