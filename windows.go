package firewall

import (
	//	"io/ioutil"
	"net"
	"strings"
)

func WindowsSetupRoutingTable(USER string) error {
	return nil
}

// Sets up specific marking for VPN-Exempt or VPN-enabled users.
func WindowsSetupNetSH(LANIP, USER string) error {
	return nil
}

func WindowsFindExePath(APP string) (string, error) {
	if out, err := Command("dir", "/s", "\"c:\\Program Files\""+APP+".exe"); err != nil {
		return "", err
	} else {
		splitprefix := strings.SplitN(out, "c:\\", 2)
		splitsuffix := strings.SplitN(splitprefix[1], APP+".exe", 2)
		return splitsuffix[0], nil
	}
}

func WindowsSetupNetSHExemptApp(LANIP, APP, LANINTERFACE, VPNINTERFACE string) error {
	if apppath, err := WindowsFindExePath(APP); err != nil {
		return err
	} else {
		if _, err := Command("netsh", "advfirewall", "firewall", "add", "rule", "name=\""+APP+"\"", "dir=in", "action=allow", "program=\""+apppath+"\"", "enable=yes", "profile=domain,private"); err != nil {
			return err
		}
		if _, err := Command("netsh", "advfirewall", "firewall", "add", "rule", "name=\""+APP+"\"", "dir=out", "action=allow", "program=\""+apppath+"\"", "enable=yes", "profile=domain,private"); err != nil {
			return err
		}
		return nil
	}
}

func WindowsSetupNetSHSecureSetup(LANIP, USER, VPNINTERFACE string) error {
	if _, err := Command("netsh", "advfirewall", "firewall", "add", "rule", "name=\""+VPNINTERFACE+"\"", "dir=in", "action=deny", "enable=yes", "profile=domain"); err != nil {
		return err
	}
	if _, err := Command("netsh", "advfirewall", "firewall", "add", "rule", "name=\""+VPNINTERFACE+"\"", "dir=out", "action=deny", "enable=yes", "profile=domain"); err != nil {
		return err
	}
	return nil
}

func WindowsCheckIPRules() bool {
	if out, err := Command("/sbin/ip", "rule", "list"); err != nil {
		return false
	} else {
		return strings.Contains(out, "0x1")
	}
}

func WindowsSetupRoutingTables(gate net.IP, USER, INTERFACE string, exempt bool, VPNINTERFACE string) error {
	GATEWAY := gate.String()
	if !exempt {
		if vpngate, err := IfIP(VPNINTERFACE); err != nil {
			return err
		} else {
			GATEWAY = vpngate.String()
		}
	}
	if !WindowsCheckIPRules() {
		if _, err := Command("route", "replace", "default", "via", GATEWAY, "table", USER); err != nil {
			return err
		}
		if _, err := Command("route", "append", "default", "via", "127.0.0.1", "dev", "lo", "table", USER); err != nil {
			return err
		}
		if _, err := Command("route", "flush", "cache"); err != nil {
			return err
		}

	}
	return nil
}

func windowsSetup(addr, gate net.IP, APP, INTERFACE, BRIF string, exempt bool, VPNINTERFACE string) error {
	LANIP := addr.String()
	if err := WindowsSetupRoutingTable(APP); err != nil {
		return err
	}
	if err := WindowsSetupNetSH(LANIP, APP); err != nil {
		return err
	}
	if exempt {
		WindowsSetupNetSHExemptApp(LANIP, APP, INTERFACE, VPNINTERFACE)
	} else {
		WindowsSetupNetSHSecureSetup(LANIP, APP, VPNINTERFACE)
	}

	if err := WindowsSetupRoutingTables(gate, APP, INTERFACE, exempt, VPNINTERFACE); err != nil {
		return err
	}
	if _, err := Command("/sbin/iptables", "-A", "OUTPUT", "!", "--src", LANIP, "-o", BRIF, "-j", "REJECT"); err != nil {
		return err
	}
	return nil
}
