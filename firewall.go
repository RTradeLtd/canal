package firewall

import (
	"fmt"
	"log"
	"runtime"
	"github.com/jackpal/gateway"
    "github.com/eyedeekay/canal/etc"
)

func Setup(user, iface string, exempt bool, vface string) error {
	if iface == "" {
		var err error
		iface, err = fwscript.DefaultIface()
		if err != nil {
			return err
		}
	}
	log.Println("firewall:", "Setting up firewall")
	gate, err := gateway.DiscoverGateway()
	if err != nil {
		return err
	}
	addr, err := fwscript.IfIP(iface)
	if err != nil {
		return err
	}
	log.Println("firewall:", iface, gate.String(), addr.String())
	switch os := runtime.GOOS; os {
	case "darwin":
	case "linux":
	case "windows":
	default:
		return fmt.Errorf("Error setting up VPN interface to be default gateway")
	}
	return nil
}
