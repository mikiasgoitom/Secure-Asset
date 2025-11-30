package valueobject

type Classification uint8

const (
	Confidential Classification = iota
	Restricted
	Internal
	Public
)

func (c Classification) String() string {
	switch c {
	case Confidential:
		return "Confidential"
	case Restricted:
		return "Restricted"
	case Internal:
		return "Internal"
	case Public:
		return "Public"
	default:
		return "Unknown"
	}
}