package dns

type Class int // 16 bit

const (
	ClassInternet  Class = 1   // RFC 1035
	ClassChaosnet  Class = 2   // RFC 1035
	ClassHesoidnet Class = 4   // RFC 1035
	ClassNone      Class = 254 // RFC 2136
)

func (class Class) String() string {

	switch class {
	case ClassInternet:
		return "Internet"
	case ClassChaosnet:
		return "Chaosnet"
	case ClassHesoidnet:
		return "Hesoidnet"
	case ClassNone:
		return ""
	default:
		return ""
	}

}
