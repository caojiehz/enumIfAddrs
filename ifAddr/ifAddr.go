package ifAddr

import (
	"github.com/wangtuanjie/ip17mon"
	"net"
)

func IsPublicIP(ip string) bool {
	IP := net.ParseIP(ip)
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 169 && ip4[1] == 254:		//保留地址
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

type LocalIfAddr struct {
	IfName string
	MTU    int
	IPv4   string
}

type PublicIP struct {
	IP string
	ip17mon.LocationInfo
}

func (addr PublicIP) Mainland() bool {
	if addr.IP == "" || !IsPublicIP(addr.IP) {
		return false
	}
	if addr.Country != string("中国") {
		return false
	}
	if addr.Region == string("香港") {
		return false
	}
	if addr.Region == string("台湾") {
		return false
	}
	if addr.Region == string("澳门") {
		return false
	}

	return true
}

type IfAddrInfo struct {
	LocalIfAddr
	PublicIP
}

func (addr IfAddrInfo) LocalIP() string {
	return addr.IPv4
}

func (addr IfAddrInfo) OuterIP() string {
	return addr.IP
}

func (addr IfAddrInfo) IsOuterIP() bool {
	if addr.IP == "" {
		return IsPublicIP(addr.LocalIfAddr.IPv4)
	}
	return IsPublicIP(addr.IP)
}
