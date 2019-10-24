package firewall

import (
	"testing"
)

func TestPlatform(t *testing.T) {
	if err := Setup("10.17.0.2", "10.17.0.1", "tun0"); err != nil {
		t.Fatal(err)
	}
}
