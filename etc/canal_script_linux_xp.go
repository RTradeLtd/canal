package fwscript

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
)

var (
	Clear = []string{
		"-F",
		"-P INPUT ACCEPT",
		"-P FORWARD ACCEPT",
		"-P OUTPUT ACCEPT",
		"-t nat -F",
		"-t mangle -F",
		"-F",
		"-X",
	}

	DenyAll = []string{
		//"-I INPUT -p tcp --dport 22 -j ACCEPT",
		"-P INPUT DROP",
		"-P FORWARD DROP",
		"-P OUTPUT DROP",
		"-A INPUT -i lo -j ACCEPT",
		"-A OUTPUT -o lo -j ACCEPT",
		"-A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT",
	}
)

func EnableVPN(tunnel, ipaddr string) []string {
	EnableVPN := []string{
		"-A INPUT  -i " + tunnel + " -j ACCEPT",
		"-A INPUT  -s " + ipaddr + " -j ACCEPT",
		"-A OUTPUT -o " + tunnel + " -j ACCEPT",
		"-A OUTPUT -d " + ipaddr + " -j ACCEPT",
	}
	return EnableVPN
}

func ExceptTCPPort(port string) []string {
	exceptPort := []string{
		"-I INPUT -p tcp --dport " + port + " -j ACCEPT",
		"-I INPUT -p tcp --dport " + port + " -j ACCEPT",
	}
	return exceptPort
}

func ExceptUDPPort(port string) []string {
	exceptPort := []string{
		"-I INPUT -p udp --dport " + port + " -j ACCEPT",
		"-I INPUT -p udp --dport " + port + " -j ACCEPT",
	}
	return exceptPort
}

func RunIPTables(net int, rules []string) (string, error) {
	nv := ""
	if net != 4 && net != 6 {
		return "", fmt.Errorf("Network must be 4 or 6")
	} else if net == 6 {
		nv = strconv.Itoa(net)
	}

	strout := ""
	for _, rule := range rules {
		strout += "ip" + nv + "tables " + rule
		output, err := Command("/sbin/ip"+nv+"tables", splitRules(rule)...)
		if err != nil {
			return strout + output, err
		}
		fmt.Printf("%s %v\n", "ip"+nv+"tables ", splitRules(rule))
	}
	return strout, nil
}

func testRunIPTables(net int, rules []string) (string, error) {
	nv := ""
	if net != 4 && net != 6 {
		return "", fmt.Errorf("Network must be 4 or 6")
	} else if net == 6 {
		nv = strconv.Itoa(net)
	}

	strout := ""
	for _, rule := range rules {
		strout += "ip" + nv + "tables " + rule
		fmt.Printf("%s %v\n", "ip"+nv+"tables ", rule)
	}
	return strout, nil
}

func SetupIPTables(rules []string) (string, error) {
	var retout string
	out1, err := RunIPTables(4, rules)
	if err != nil {
		return "", err
	}
	out2, err := RunIPTables(6, rules)
	if err != nil {
		return "", err
	}
	retout += out1
	retout += out2
	return retout, nil
}

func ResetIPTables() (string, error) {
	return SetupIPTables(Clear)
}

func DenyIPTables() (string, error) {
	return SetupIPTables(DenyAll)
}

func EnableVPNIPTables(tunnel, ipaddr string) (string, error) {
	return SetupIPTables(EnableVPN(tunnel, ipaddr))
}

func RunExceptTCPPort(port string) (string, error) {
	return SetupIPTables(ExceptTCPPort(port))
}

func RunExceptUDPPort(port string) (string, error) {
	return SetupIPTables(ExceptUDPPort(port))
}

func LinuxSetupVPNTunnelled(tunnel, ipaddr string) (string, error) {
	var retout string
	out1, err := ResetIPTables()
	if err != nil {
		return "", err
	}
	out2, err := DenyIPTables()
	if err != nil {
		return "", err
	}
	out3, err := EnableVPNIPTables(tunnel, ipaddr)
	if err != nil {
		return "", err
	}
	retout += out1
	retout += out2
	retout += out3
	return retout, nil
}

func LinuxSetupVPNExceptPort(tunnel, ipaddr, port string) (string, error) {
	var retout string
	out1, err := LinuxSetupVPNTunnelled(tunnel, ipaddr)
	if err != nil {
		return "", err
	}
	out2, err := RunExceptUDPPort(port)
	if err != nil {
		return "", err
	}
	out3, err := RunExceptTCPPort(port)
	if err != nil {
		return "", err
	}
	retout += out1
	retout += out2
	retout += out3
	return retout, nil
}

func LinuxSetupVPNDNS(tunnel, ipaddr string) (string, error) {
	return LinuxSetupVPNExceptPort(tunnel, ipaddr, "53")
}

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
