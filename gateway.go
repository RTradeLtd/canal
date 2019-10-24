package firewall

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
)

func LinuxServerSetup(iface, gface string) error {
	if err := ioutil.WriteFile("/proc/sys/net/ipv4/ip_forward", []byte("1"), 0644); err != nil {

	}
	if err := exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", gface, "-j", "MASQUERADE").Run(); err != nil {
		return err
	}
	if err := exec.Command("iptables", "-A", "FORWARD", "-i", iface, "-o", gface, "-j", "ACCEPT").Run(); err != nil {
		return err
	}
	if err := exec.Command("iptables", "-A", "FORWARD", "-i", gface, "-o", iface, "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT").Run(); err != nil {
		return err
	}
	return nil
}

func WindowsServerSetup(iface, gface string) error {
	//if err := exec.Command("", "").Run(); err != nil {
        //return err
	//}
	return fmt.Errorf("%s is not supported yet as a server, please check back soon", "Windows")
}

func DarwinServerSetup(iface, gface string) error {
	//if err := exec.Command("", "").Run(); err != nil {
        //return err
	//}
	return fmt.Errorf("%s is not supported yet as a server, please check back soon", "OSX")
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
