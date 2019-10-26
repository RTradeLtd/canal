package firewall

import (
	"fmt"
	"net"
	"runtime"
	"strings"

	"github.com/jackpal/gateway"
)

func windowsSetup(addr, gate net.IP, user, iface string, exempt bool, vface string) error {
	return nil
}

func darwinSetup(addr, gate net.IP, user, iface string, exempt bool, vface string) error {
	return nil
}

func IfIP(INTERFACE string) (net.IP, error) {
	if i, err := net.InterfaceByName(INTERFACE); err == nil {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if !strings.HasPrefix(ip.String(), "127.0.") {
				return ip, nil
			}
		}
	}
	return nil, fmt.Errorf("IP Address for interface not found")
}

func Setup(user, iface string, exempt bool, vface string) error {
	gate, err := gateway.DiscoverGateway()
	if err != nil {
		return err
	}
	addr, err := IfIP(iface)
	if err != nil {
		return err
	}
	switch os := runtime.GOOS; os {
	case "darwin":
		return darwinSetup(addr, gate, user, iface, exempt, vface)
	case "linux":
		return linuxSetup(addr, gate, user, iface, "br0", exempt, vface)
	case "windows":
		return windowsSetup(addr, gate, user, iface, exempt, vface)
	default:
		return fmt.Errorf("Error setting up VPN interface to be default gateway")
	}
	return nil
}
