package fwscript

import (
	"fmt"
	"strconv"
)

var AllowBasic = []string{
	"advfirewall firewall add rule name=\"Core Networking (DNS-Out)\" dir=out action=allow protocol=UDP remoteport=53 program=\"%%%%systemroot%%%%\\system32\\svchost.exe\" service=\"dnscache\"",
	"advfirewall firewall add rule name=\"Core Networking (DHCP-Out)\" dir=out action=allow protocol=UDP localport=68 remoteport=67 program=\"%%%%systemroot%%%%\\system32\\svchost.exe\" service=\"dhcp\"",
}

func VPNRouteWithMetric(gateway string, metric int, ifname string) []string {
	return []string{
		"ADD 0.0.0.0 MASK 0.0.0.0 " + gateway + " METRIC " + strconv.Itoa(metric) + " IF " + ifname,
	}
}

func VPNRoute(gateway string, ifname string) []string {
	return VPNRouteWithMetric(gateway, 1, ifname)
}

func ExceptWindowsTCPPort(port int) []string {
	return []string{
		"advfirewall firewall add rule name=\"allow tcp port " + strconv.Itoa(port) + " in\" dir=in action=allow protocol=TCP localport=" + strconv.Itoa(port),
		"advfirewall firewall add rule name=\"allow tcp port " + strconv.Itoa(port) + " out\" dir=in action=allow protocol=TCP localport=" + strconv.Itoa(port),
	}
}

func ExceptWindowsUDPPort(port int) []string {
	return []string{
		"advfirewall firewall add rule name=\"allow udp port " + strconv.Itoa(port) + " in\" dir=in action=allow protocol=UDP localport=" + strconv.Itoa(port),
		"advfirewall firewall add rule name=\"allow udp port " + strconv.Itoa(port) + " out\" dir=in action=allow protocol=UDP localport=" + strconv.Itoa(port),
	}
}

func ExceptApplication(appname, pathtoapp string) []string {
	return []string{
		"advfirewall firewall add rule name=\"allow " + appname + "\" dir=in program=\"" + pathtoapp + "\" security=authnoencap action=allow",
	}
}

func RunRoute(rules []string) (string, error) {
	strout := ""
	for _, rule := range rules {
		strout += "route.exe"
		output, err := Command("route.exe", splitRules(rule)...)
		if err != nil {
			return strout + output, err
		}
		fmt.Printf("%s %v\n", "route.exe ", splitRules(rule))
	}
	return strout, nil
}

func testRunRoute(rules []string) (string, error) {
	strout := ""
	for _, rule := range rules {
		strout += "netsh.exe" + rule
		fmt.Printf("%s %v\n", "netsh.exe", rule)
	}
	return strout, nil
}

func RunNetSH(rules []string) (string, error) {
	strout := ""
	for _, rule := range rules {
		strout += "netsh.exe"
		output, err := Command("netsh.exe", splitRules(rule)...)
		if err != nil {
			return strout + output, err
		}
		fmt.Printf("%s %v\n", "netsh.exe ", splitRules(rule))
	}
	return strout, nil
}

func testRunNetSH(rules []string) (string, error) {
	strout := ""
	for _, rule := range rules {
		strout += "netsh.exe" + rule
		fmt.Printf("%s %v\n", "netsh.exe", rule)
	}
	return strout, nil
}

func RunExceptDHCP() (string, error) {
	return RunNetSH(AllowBasic)
}

func RunExceptWindowsTCPPort(port int) (string, error) {
	return RunNetSH(ExceptWindowsTCPPort(port))
}

func RunExceptWindowsUDPPort(port int) (string, error) {
	return RunNetSH(ExceptWindowsUDPPort(port))
}

func RunExceptWindowsApplication(appname, pathtoapp string) (string, error) {
	return RunNetSH(ExceptApplication(appname, pathtoapp))
}

func SetupRoute(gateway string, ifname string) (string, error) {
	return RunRoute(VPNRoute(gateway, ifname))
}

func SetupRouteWithMetric(gateway string, metric int, ifname string) (string, error) {
	return RunRoute(VPNRouteWithMetric(gateway, metric, ifname))
}

func WindowsSetupRouteDHCP(gateway string, ifname string) (string, error) {
	strout := ""
	out1, err := SetupRoute(gateway, ifname)
	if err != nil {
		return "", err
	}
	out2, err := RunExceptDHCP()
	if err != nil {
		return "", err
	}
	out3, err := RunExceptWindowsUDPPort(53)
	if err != nil {
		return "", err
	}
	strout += out1
	strout += out2
	strout += out3
	return strout, nil
}

func WindowsServerSetup(iface, gface string) error {
	/*if err := exec.Command("netsh", "interface", "ipv4", "set", "interface", "name="+gface, "forwarding=enabled").Run(); err != nil {
	        return err
		}*/
	return fmt.Errorf("%s is not supported yet as a server, please check back soon", "Windows")
}
