package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/demtoni/tade/internal/manager"
)

var (
	// manager related flags
	flagManager  = flag.String("manager", "", "manager address (required)")
	flagSecret   = flag.String("secret", "", "server secret to protect api (required)")
	flagHostname = flag.String("hostname", "", "public hostname/ip of the machine to generate config URIs, defaults to server address")
	flagState    = flag.String("state", "", "path to state file (required)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [flags]\nwhere flags are:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()

	if *flagManager == "" || *flagSecret == "" || *flagState == "" {
		usage()
	}

	manager.Addr = *flagManager
	manager.Secret = *flagSecret
	manager.PathToState = *flagState

	var err error

	if *flagHostname == "" {
		if *flagHostname, _, err = net.SplitHostPort(*flagManager); err != nil {
			log.Fatal(err)
		}
	}

	manager.Hostname = *flagHostname

	m, err := manager.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(m.Serve())
}
