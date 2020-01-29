package fwscript

import (
	"testing"
)

func TestWindowsPsList(t *testing.T) {
	pid, err := GetPidString("i2p")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(pid test) ", pid)
	route, err := testRunRoute(VPNRouteWithMetric("10.17.0.2", 1, "tun"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(route test) ", route)
	tdns, err := testRunNetSH(ExceptWindowsTCPPort(53))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(dns test) ", tdns)
	udns, err := testRunNetSH(ExceptWindowsUDPPort(53))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(dns stest)", udns)
	dhcp, err := testRunNetSH(AllowBasic)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(dhcp test) ", dhcp)
}
