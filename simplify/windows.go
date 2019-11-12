package firewall

import (
	//	"io/ioutil"
	"log"
	//	"math/rand"
	"net"
	//	"strconv"

	"github.com/jackpal/gateway"
)

func WindowsSetupRoutingMetrics(ip, gate net.IP) error {
	//GATEWAY := gate.String()
	IP := ip.String()
	if _, err := Command("route.exe", "flush", "cache"); err != nil {
		return err
	}
	if _, err := Command("route.exe", "change", "0.0.0.0", "mask", "0.0.0.0", IP, "metric", "5"); err != nil {
		return err
	}
	return nil
}

func WindowsSetupRoutingTable(vpnaddr, vpngate net.IP) error {
	VPNIP := vpnaddr.String()
	VPNGATE := vpngate.String()
	if _, err := Command("route.exe", "add", VPNIP, VPNGATE, "metric", "1", "1.1.1.1", "route", "add", "default", "gw", "10.0.0.2", "metric 2"); err != nil {
		return err
	}
	return nil
}

func windowsSetup(addr, vpnaddr, gate, vpngate net.IP) error {
	if err := WindowsSetupRoutingMetrics(addr, gate); err != nil {
		return err
	}
	if err := WindowsSetupRoutingTable(vpnaddr, vpngate); err != nil {
		return err
	}
	return nil
}

func SetupWindows(vpnaddr, vpngate net.IP) error {
	iface, err := DefaultIface()
	if err != nil {
		return err
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
	return windowsSetup(addr, vpnaddr, gate, vpngate)
}
