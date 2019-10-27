package firewall

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jackpal/gateway"
)

func darwinSetup(addr, gate net.IP, user, iface, bridge string, exempt bool, vface string) error {
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

func Command(command string, args ...string) (string, error) {
	log.Println(command, args)
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func DefaultIface() (string, error) {
	gate, err := gateway.DiscoverGateway()
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.String() == gate.String() {
				return i.Name, nil
			}
			// process IP address
		}
	}
	return "", nil
}

func Setup(user, iface string, exempt bool, vface string) error {
	if iface == "" {
		var err error
		iface, err = DefaultIface()
		if err != nil {
			return err
		}
	}
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
		return darwinSetup(addr, gate, user, iface, "bridge", exempt, vface)
	case "linux":
		return linuxSetup(addr, gate, user, iface, "br0", exempt, vface)
	case "windows":
		return windowsSetup(addr, gate, user, iface, "bridge", exempt, vface)
	default:
		return fmt.Errorf("Error setting up VPN interface to be default gateway")
	}
	return nil
}
