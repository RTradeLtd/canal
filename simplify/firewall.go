package firewall

import (
	//	"io/ioutil"
	"log"
	//	"math/rand"
	"fmt"
	"net"
	"os/exec"
	//	"strconv"
    "runtime"
	"strings"

	"github.com/jackpal/gateway"
)

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
		log.Println("firewall: Error", command, "\n\t", args, "\n\t", string(output), "\n\t", err)
		return "", fmt.Errorf("firewall: Error %s %s", output, err)
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

func Setup(vpnaddr, vpngate net.IP) error {
	switch os := runtime.GOOS; os {
	case "darwin":
		//return darwinSetup(addr, gate, user, iface, "bridge", exempt, vface)
        return fmt.Errorf("%s", "Mac support isn't ready yet")
	case "linux":
		return SetupLinux(vpnaddr, vpngate)
	case "windows":
		return SetupWindows(vpnaddr, vpngate)
	default:
		return fmt.Errorf("Error setting up VPN interface to be default gateway")
	}
	return nil
}

