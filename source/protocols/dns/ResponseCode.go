package dns

type ResponseCode int // 4 bit

const (
	ResponseCodeNoError                       = 0 // RFC 1035, 4.1.1
	ResponseCodeFormatError                   = 1 // RFC 1035, 4.1.1
	ResponseCodeServerFailure                 = 2 // RFC 1035, 4.1.1
	ResponseCodeNonExistDomain                = 3 // RFC 1035, 4.1.1
	ResponseCodeNotImplemented                = 4 // RFC 1035, 4.1.1
	ResponseCodeQueryRefused                  = 5 // RFC 1035, 4.1.1
	ResponseCodeWhyExistDomain                = 6  // RFC 2136
	ResponseCodeWhyExistRecord                = 7  // RFC 2136
	ResponseCodeNotExistRecord                = 8  // RFC 2136
	ResponseCodeNotAuthorized                 = 9  // RFC 2136 and RFC 8945
	ResponseCodeNotZone                       = 10 // RFC 2136
	ResponseCodeServerOperationNotImplemented = 11 // RFC 8490
	ResponseCodeBadVersion                    = 16 // RFC 6891 and RFC 8945
	ResponseCodeBadSignature                  = 17 // RFC 8945
	ResponseCodeBadTimeWindow                 = 18 // RFC 8945
	ResponseCodeBadKeyMode                    = 19 // RFC 2930
	ResponseCodeBadKeyName                    = 20 // RFC 2930
	ResponseCodeBadAlgorithm                  = 21 // RFC 2930
	ResponseCodeBadTruncation                 = 22 // RFC 8945
	ResponseCodeBadCookie                     = 23 // RFC 7873
)

func (rcode ResponseCode) String() string {

	switch rcode {
	case ResponseCodeNoError:
		return "No Error"
	case ResponseCodeFormatError:
		return "Format Error"
	case ResponseCodeServerFailure:
		return "Server Failure"
	case ResponseCodeNonExistDomain:
		return "Non-Existent Domain"
	case ResponseCodeNotImplemented:
		return "Not Implemented"
	case ResponseCodeQueryRefused:
		return "Query Refused"
	case ResponseCodeWhyExistDomain:
		return "Domain Exists When It Should Not"
	case ResponseCodeWhyExistRecord:
		return "Record Exists When It Should Not"
	case ResponseCodeNotExistRecord:
		return "Non-Existent Record Set"
	case ResponseCodeNotAuthorized:
		return "Not Authorized"
	case ResponseCodeNotZone:
		return "Domain Not Contained In Zone"
	case ResponseCodeServerOperationNotImplemented:
		return "Server Operation Not Implemented"
	case ResponseCodeBadVersion:
		return "Bad Version"
	case ResponseCodeBadSignature:
		return "Bad Signature"
	case ResponseCodeBadTimeWindow:
		return "Bad Time Window"
	case ResponseCodeBadKeyMode:
		return "Bad Key Mode"
	case ResponseCodeBadKeyName:
		return "Bad Key Name"
	case ResponseCodeBadAlgorithm:
		return "Bad Algorithm"
	case ResponseCodeBadTruncation:
		return "Bad Truncation"
	case ResponseCodeBadCookie:
		return "Bad Cookie"
	default:
		return ""
	}

}
