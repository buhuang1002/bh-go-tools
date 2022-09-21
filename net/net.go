package bh_net

import (
	"net"
)

func GetPrivateIp() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []string{}
	}
	var ipList []string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && IsPrivateIp(ipnet.IP) {
				ipList = append(ipList, ipnet.IP.String())
			}
		}
	}
	return ipList
}

func IsPrivateIp(IP net.IP) bool {
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return false
}
