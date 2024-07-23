package net

const (
	NETWORK_UNKNOWN uint32 = iota
	NETWORK_TCP
	NETWORK_TCP4
	NETWORK_TCP6
)

func ToNetworkTypeString(networkType uint32) string {
	switch networkType {
	case NETWORK_TCP:
		return "tcp"
	case NETWORK_TCP4:
		return "tcp4"
	case NETWORK_TCP6:
		return "tcp6"
	default:
		panic("unsupported network type")
	}
}

func ToNetworkTypeInt(networkType string) uint32 {
	switch networkType {
	case "tcp":
		return NETWORK_TCP
	case "tcp4":
		return NETWORK_TCP4
	case "tcp6":
		return NETWORK_TCP6
	default:
		panic("unsupported network type")
	}
}
