package fwscript

import (
	"testing"
)

func TestPsList(t *testing.T) {
	pid, err := GetPidString("i2p")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(pid test) ", pid)
	clear4, err := testRunIPTables(4, Clear)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(clear4 test) ", clear4)
	clear6, err := testRunIPTables(6, Clear)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(clear6 test) ", clear6)
	deny4, err := testRunIPTables(4, DenyAll)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(deny4 test) ", deny4)
	deny6, err := testRunIPTables(6, DenyAll)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(deny6 test) ", deny6)
	enable4, err := testRunIPTables(4, EnableVPN("tun0", "10.17.0.2"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("(vpn4 test) ", enable4)
}
