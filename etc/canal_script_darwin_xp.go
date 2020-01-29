package fwscript

import (
	"fmt"
)

func DarwinServerSetup(iface, gface string) error {
	/*if err := exec.Command("sysctl", "-w", "net.inet.ip.forwarding=1").Run(); err != nil {
	        return err
		}*/
	return fmt.Errorf("%s is not supported yet as a server, please check back soon", "OSX")
}
