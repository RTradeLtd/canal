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
	log.Println("firewall: looking up IP of", INTERFACE)
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
	} else {
		return nil, err
	}
	return nil, fmt.Errorf("Undefined error discovering IP of", INTERFACE)
}

func Command(command string, args ...string) (string, error) {
	log.Println("firewall:", command, args)
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
			log.Println("Does this err matter?")
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
			h := strings.Split(ip.String(), ".")
			g := strings.Split(gate.String(), ".")

			if len(h) == 4 && len(g) == 4 {
				i2 := h[0] + h[1] + h[2]
				g2 := g[0] + g[1] + g[2]
				if i2 == g2 {
					log.Println("firewall:", i2, "==", g2)
					return i.Name, nil
				}
				log.Println("firewall:", i2, "!=", g2)
			}
			// process IP address
		}
	}
	return "", fmt.Errorf("default interface not found")
}

func Setup(user, iface string, exempt bool, vface string) error {
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
