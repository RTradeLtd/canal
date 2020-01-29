package fwscript

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Setup(iface string) error {
	if iface == "" {
		var err error
		iface, err = DefaultIface()
		if err != nil {
			return err
		}
	}
	addr, err := IfIP(iface)
	if err != nil {
		return err
	}
	gate := tunGate(addr.String())
	log.Println("firewall:", iface, gate, addr.String())
	switch os := runtime.GOOS; os {
	case "darwin":
		return fmt.Errorf("Unsupported platform error, come back later.")
	case "linux":
		_, err := LinuxSetupVPNDNS(iface, addr.String())
		return err
	case "windows":
		_, err := WindowsSetupRouteDHCP(gate, addr.String())
		return err
	default:
		return fmt.Errorf("Error setting up VPN interface to be default gateway")
	}
	return nil
}

func ServerSetup(iface, gface string) error {
	switch os := runtime.GOOS; os {
	case "darwin":
		return DarwinServerSetup(iface, gface)
	case "linux":
		return LinuxServerSetup(iface, gface)
	case "windows":
		return WindowsServerSetup(iface, gface)
	default:
		return fmt.Errorf("")
	}
	return nil
}

func tunGate(addr string) string {
	dec := strings.Split(addr, ".")
	if len(dec) == 4 {
		return dec[0] + "." + dec[1] + "." + dec[2] + "." + "1"
	}
	return "192.168.0.1"
}
