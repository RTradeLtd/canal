package main

import (
	"flag"
	"log"

	"github.com/eyedeekay/canal/etc"
)

var (
	tunnelname = flag.String("tun", "tun0", "Name of the tunnel to auto-configure")
	gateway    = flag.String("gate", fwscript.DefaultGate(), "Gateway to forward traffic to")
	server     = flag.Bool("server", false, "Configure a server (default false)")
)

func main() {
	flag.Parse()
	if *server {
		if err := fwscript.ServerSetup(*tunnelname, *gateway); err != nil {
			log.Fatal(err)
		}
		log.Println("VPN Server Setup")
	} else {
		if err := fwscript.Setup(*tunnelname); err != nil {
			log.Fatal(err)
		}
		log.Println("VPN Client Setup")
	}
}
