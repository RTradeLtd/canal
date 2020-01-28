package fwscript

import (
    "fmt"
	"strconv"
)

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

func RunNetSH(net int, rules []string) (string, error) {
	nv := ""
	if net != 4 && net != 6 {
		return "", fmt.Errorf("Network must be 4 or 6")
	} else if net == 6 {
		nv = strconv.Itoa(net)
	}

	strout := ""
	for _, rule := range rules {
		strout += "ip" + nv + "tables " + rule
		output, err := Command("netsh.exe", splitRules(rule)...)
		if err != nil {
			return strout + output, err
		}
		fmt.Printf("%s %v\n", "ip"+nv+"tables ", splitRules(rule))
	}
	return strout, nil
}

func testRunNetSH(net int, rules []string) (string, error) {
	nv := ""
	if net != 4 && net != 6 {
		return "", fmt.Errorf("Network must be 4 or 6")
	} else if net == 6 {
		nv = strconv.Itoa(net)
	}

	strout := ""
	for _, rule := range rules {
		strout += "ip" + nv + "tables " + rule
		fmt.Printf("%s %v\n", "netsh.exe", rule)
	}
	return strout, nil
}