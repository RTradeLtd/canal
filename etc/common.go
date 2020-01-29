package fwscript

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jackpal/gateway"
)

func splitRules(rule string) []string {
	return strings.Split(rule, " ")
}

func GetPidString(Executable string) (string, error) {
	id, err := GetPidOf(Executable)
	return strconv.Itoa(id), err
}

func GetPidOf(Executable string) (int, error) {
	processes, err := ps.Processes()
	if err != nil {
		return 0, err
	}
	for _, proc := range processes {
		if strings.Contains(proc.Executable(), Executable) {
			return proc.Pid(), nil
		}
	}
	return 0, fmt.Errorf("(setup) error %s is not running", Executable)
}

func AppendFile(filename, text string, perms os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, perms)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		return err
	}
	return nil
}

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
	return nil, fmt.Errorf("Undefined error discovering IP of %s", INTERFACE)
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

func DefaultGate() string {
	gate, _ := gateway.DiscoverGateway()
	return gate.String()
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
