package firewall

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

func WindowsSetup(addr, gate, iface string) error {
	if err := exec.Command("route", "-f").Run(); err != nil {
		return err
	}
	if interfaces, err := net.Interfaces(); err != nil {
		for _, val := range interfaces {
			if val.Name != iface {
				if val.Flags == net.FlagUp && val.Flags != net.FlagLoopback {
					if err := exec.Command("route", "ADD", "0.0.0.0", "MASK", "0.0.0.0", gate, iface).Run(); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func LinuxSetup(addr, gate, iface string) error {
	if err := exec.Command("/sbin/ip", "route", "del", "default").Run(); err != nil {
		return err
	}
	if err := exec.Command("/sbin/ip", "route", "add", "default", "via", gate, "dev", iface).Run(); err != nil {
		return err
	}
	if interfaces, err := net.Interfaces(); err != nil {
		for _, val := range interfaces {
			if val.Name != iface {
				if val.Flags == net.FlagUp && val.Flags != net.FlagLoopback {
					if err := exec.Command("/sbin/ip", "route", "add", addr, "via", gate, "dev", val.Name).Run(); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func DarwinSetup(addr, gate, iface string) error {
	if err := exec.Command("route", "-n", "flush").Run(); err != nil {
		return err
	}
	if err := exec.Command("route", "add", "-ifscope", iface, addr, gate).Run(); err != nil {
		return err
	}
	if interfaces, err := net.Interfaces(); err != nil {
		for _, val := range interfaces {
			if val.Name != iface {
				if val.Flags == net.FlagUp && val.Flags != net.FlagLoopback {
					if err := exec.Command("route", "add", "-ifscope", val.Name, addr, gate, "0.0.0.0").Run(); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func Setup(addr, gate, iface string) error {
	switch os := runtime.GOOS; os {
	case "darwin":
		return DarwinSetup(addr, gate, iface)
	case "linux":
		return LinuxSetup(addr, gate, iface)
	case "windows":
		return WindowsSetup(addr, gate, iface)
	default:
		return fmt.Errorf("Error setting up VPN interface to be default gateway")
	}
	return nil
}
