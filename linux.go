package firewall

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

func LinuxFindNum() string {
	if bytes, err := ioutil.ReadFile("/etc/iproute2/rt_tables"); err == nil {
		for {
			value := strconv.Itoa(rand.Intn(255))
			if !strings.Contains(string(bytes), value) {
				return value
			}
		}
	}
	return "199"
}

func LinuxFindNext() string {
	var retv int
	if bytes, err := ioutil.ReadFile("/etc/iproute2/rt_tables"); err == nil {
		lines := strings.Split(string(bytes), "\n")
		for _, content := range lines {
			if !strings.HasPrefix(content, "#") {
				out := strings.Split(content, "    ")
				log.Println("firewall: finding routing table", out, len(out))
				if len(out) > 2 {
					log.Println("firewall: table", out[2])

					retv, _ = strconv.Atoi(strings.Trim(out[2], " "))
				}
			}
		}
	}
	return strconv.Itoa(retv + 1)
}

func LinuxSetupRoutingTable(USER string) error {
	if bytes, err := ioutil.ReadFile("/etc/iproute2/rt_tables"); err == nil {
		if !strings.Contains(string(bytes), USER) {

			if err := AppendFile("/etc/iproute2/rt_tables", "\n"+LinuxFindNum()+"    "+USER+"    "+LinuxFindNext(), 0644); err != nil {
				return err
			}
		}
	} else {
		return err
	}
	return nil
}

// Sets up specific marking for VPN-Exempt or VPN-enabled users.
func LinuxFlushIPTables() error {
	if _, err := Command("/sbin/iptables", "-F", "-t", "nat"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-F", "-t", "mangle"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-F", "-t", "filter"); err != nil {
		return err
	}
	return nil
}

// Sets up specific marking for VPN-Exempt or VPN-enabled users.
func LinuxSetupIPTables(LANIP, USER string) error {
	if _, err := Command("/sbin/iptables", "-t", "mangle", "-A", "OUTPUT", "!", "--dest", LANIP, "-m", "owner", "--uid-owner", USER, "-j", "MARK", "--set-mark", "0x1"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-t", "mangle", "-A", "OUTPUT", "--dest", LANIP, "-p", "udp", "--dport", "53", "-m", "owner", "--uid-owner", USER, "-j", "MARK", "--set-mark", "0x1"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-t", "mangle", "-A", "OUTPUT", "--dest", LANIP, "-p", "tcp", "--dport", "53", "-m", "owner", "--uid-owner", USER, "-j", "MARK", "--set-mark", "0x1"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-t", "mangle", "-A", "OUTPUT", "!", "--src", LANIP, "-j", "MARK", "--set-mark", "0x1"); err != nil {
		return err
	}
	return nil
}

func LinuxSetupIPTablesExemptUser(LANIP, USER, LANINTERFACE, VPNINTERFACE string) error {
	if _, err := Command("/sbin/iptables", "-A", "OUTPUT", "-o", "lo", "-m", "owner", "--uid-owner", USER, "-j", "ACCEPT"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-A", "OUTPUT", "-o", VPNINTERFACE, "-m", "owner", "--uid-owner", USER, "-j", "ACCEPT"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-A", "OUTPUT", "-o", LANINTERFACE, "-m", "owner", "--uid-owner", USER, "-j", "ACCEPT"); err != nil {
		return err
	}
	return nil
}

func LinuxSetupIPTablesTagUser(LANIP, USER, VPNINTERFACE string) error {
	//# allow bittorrent
	if _, err := Command("/sbin/iptables", "-A", "INPUT", "-i", VPNINTERFACE, "-p", "tcp", "--dport", "59560", "-j", "ACCEPT"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-A", "INPUT", "-i", VPNINTERFACE, "-p", "tcp", "--dport", "6443", "-j", "ACCEPT"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-A", "INPUT", "-i", VPNINTERFACE, "-p", "udp", "--dport", "8881", "-j", "ACCEPT"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-A", "INPUT", "-i", VPNINTERFACE, "-p", "udp", "--dport", "7881", "-j", "ACCEPT"); err != nil {
		return err
	}
	//# send DNS to quadnine or cloudflare for $VPNUSER
	if _, err := Command("/sbin/iptables", "-t", "nat", "-A", "OUTPUT", "--dest", LANIP, "-p", "udp", "--dport", "53", "-m", "owner", "--uid-owner", USER, "-j", "DNAT", "--to-destination", "9.9.9.9"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-t", "nat", "-A", "OUTPUT", "--dest", LANIP, "-p", "tcp", "--dport", "53", "-m", "owner", "--uid-owner", USER, "-j", "DNAT", "--to-destination", "1.1.1.1"); err != nil {
		return err
	}
	// Allow user access to lo and VPNINTERFACE
	if _, err := Command("/sbin/iptables", "-A", "OUTPUT", "-o", "lo", "-m", "owner", "--uid-owner", USER, "-j", "ACCEPT"); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-A", "OUTPUT", "-o", VPNINTERFACE, "-m", "owner", "--uid-owner", USER, "-j", "ACCEPT"); err != nil {
		return err
	}
	return nil
}

func LinuxCheckIPRules() bool {
	if out, err := Command("/sbin/ip", "rule", "list"); err != nil {
		return false
	} else {
		return strings.Contains(out, "0x1")
	}
}

func LinuxSetupRoutingTables(gate net.IP, USER, INTERFACE string, exempt bool, VPNINTERFACE string) error {
	GATEWAY := gate.String()
	if !exempt {
		if vpngate, err := IfIP(VPNINTERFACE); err != nil {
			return err
		} else {
			if _, err := Command("/sbin/ip", "route", "replace", "default", "via", GATEWAY); err != nil {
				return err
			}
			GATEWAY = vpngate.String()
		}
	}
	if !LinuxCheckIPRules() {
		if _, err := Command("/sbin/ip", "route", "replace", "default", "via", GATEWAY, "table", USER); err != nil {
			return err
		}
		if _, err := Command("/sbin/ip", "route", "append", "default", "via", "127.0.0.1", "dev", "lo", "table", USER); err != nil {
			return err
		}
		if _, err := Command("/sbin/ip", "route", "flush", "cache"); err != nil {
			return err
		}

	}
	return nil
}

func linuxSetup(addr, gate net.IP, USER, INTERFACE, BRIF string, exempt bool, VPNINTERFACE string) error {
	LANIP := addr.String()
	if err := LinuxSetupRoutingTable(USER); err != nil {
		return err
	}
	/*if err := LinuxFlushIPTables(); err != nil {
		return err
	}*/
	if err := LinuxSetupIPTables(LANIP, USER); err != nil {
		return err
	}
	if exempt {
		LinuxSetupIPTablesExemptUser(LANIP, USER, INTERFACE, VPNINTERFACE)
	} else {
		LinuxSetupIPTablesTagUser(LANIP, USER, VPNINTERFACE)
	}

	if err := LinuxSetupRoutingTables(gate, USER, INTERFACE, exempt, VPNINTERFACE); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-A", "OUTPUT", "!", "--src", LANIP, "-o", BRIF, "-j", "REJECT"); err != nil {
		return err
	}
	return nil
}
