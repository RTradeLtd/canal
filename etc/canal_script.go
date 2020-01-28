package fwscript

import (
	"fmt"
	"github.com/jackpal/gateway"
	"log"
	"runtime"
)

func Setup(iface string) error {
	if iface == "" {
		var err error
		iface, err = DefaultIface()
		if err != nil {
			return err
		}
	}
	log.Println("firewall:", "Setting up firewall")
	gate, err := gateway.DiscoverGateway()
	if err != nil {
		return err
	}
	addr, err := IfIP(iface)
	if err != nil {
		return err
	}
	log.Println("firewall:", iface, gate.String(), addr.String())
	switch os := runtime.GOOS; os {
	case "darwin":
	case "linux":
		LinuxSetupVPNDNS(iface, addr.String())
	case "windows":
		//WindowsSetupVPNDNS(iface, addr.String())
	default:
		return fmt.Errorf("Error setting up VPN interface to be default gateway")
	}
	return nil
}
