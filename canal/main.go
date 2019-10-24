package main

import (
	"flag"
	"github.com/eyedeekay/canal"
)

var (
	gate     = flag.String("gatewaytunnel", "tun0", "")
	tun      = flag.String("servertunnel", "tun1", "")
	gateaddr = flag.String("gatewayaddress", "10.79.0.1", "")
	maskaddr = flag.String("maskaddress", "0.0.0.0", "")
)

func main() {
	flag.Parse()
	if err := firewall.Setup(*maskaddr, *gateaddr, *tun); err != nil {
		panic(err)
	}
	if err := firewall.ServerSetup(*tun, *gate); err != nil {
		panic(err)
	}
}
