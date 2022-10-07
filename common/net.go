package common

import (
	"bytes"
	"errors"
	"net"
	"os"
)

func GetMacAddr() string {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				return i.HardwareAddr.String()
			}
		}
	}
	return ""
}

// GetIPAddr returns the public IP address of this machine (if any), returns
// private IP otherwise.
func GetIPAddr() string {
	ipList := getIPList()

	// Try to return public IP first
	for _, ip := range ipList {
		private, err := isPrivateIP(ip)
		if err != nil || private {
			continue
		}

		return ip
	}

	// Return any IP if no public IP was found
	for _, ip := range ipList {
		return ip
	}

	return ""
}

// IsMyIP checks if the specified IP is a known IP address of this machine.
func IsMyIP(ip string) bool {
	ipList := getIPList()
	for _, myIP := range ipList {
		if ip == myIP {
			return true
		}
	}

	return false
}

func getIPList() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	ipList := make([]string, 0, 8)
	for _, itf := range interfaces {
		if itf.Name != "lo" {
			addrs, err := itf.Addrs()
			if err != nil {
				os.Stderr.WriteString("Oops: " + err.Error() + "\n")
				os.Exit(1)
			}

			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ipList = append(ipList, ipnet.IP.String())
					}
				}
			}
		}
	}

	return ipList
}

func isPrivateIP(ip string) (bool, error) {
	private := false

	IP := net.ParseIP(ip)
	if IP == nil {
		return private, errors.New("Not a valid IP")
	}

	_, private24BitBlock, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		return private, err
	}

	_, private20BitBlock, err := net.ParseCIDR("172.16.0.0/12")
	if err != nil {
		return private, err
	}

	_, private16BitBlock, err := net.ParseCIDR("192.168.0.0/16")
	if err != nil {
		return private, err
	}

	private = private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)

	return private, nil
}
