package handlerutils

import (
	"bytes"
	"net"
	"net/http"
	"strings"
)

type IPRange struct {
	start net.IP
	end   net.IP
}

func IsIPInRange(r IPRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var LanIPRanges = []IPRange{
	IPRange{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	IPRange{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	IPRange{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
}

// IsLanIpAddress - check to see if this ip is in a private subnet
func IsLanIpAddress(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range LanIPRanges {
			// check if this ip is in a private range
			if IsIPInRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func GetClientIPAddress(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || IsLanIpAddress(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return ""
}
