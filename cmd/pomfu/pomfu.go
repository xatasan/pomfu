package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xatasan/pomfu"
)

var (
	noConf, listSrv, html bool
	server                string
	srv                   *pomfu.Pomf
)

// start pomfu: parse flags and either list arguments (if -l was
// specified) or start uploading.
func main() {
	flag.BoolVar(&noConf, "n", false, "don't read the configuration file")
	flag.BoolVar(&listSrv, "l", false, "list available servers")
	flag.BoolVar(&html, "H", false, "require random server to support html")
	flag.StringVar(&server, "s", "", "use this server to upload")
	flag.Parse()

	if !noConf {
		pomfu.ReadConfig()
	}

	if server != "" {
		var ok bool
		srv, ok = pomfu.Servers[server]
		if !ok {
			fmt.Fprintf(os.Stderr, "No server with key: \"%s\"\n", server)
			os.Exit(1)
		}
	}

	switch {
	case listSrv:
		list()
	default:
		upload(flag.Args())
	}
}
