package main

import "flag"

var (
	noConf, listSrv, html bool
	server                string
)

func main() {
	flag.BoolVar(&noConf, "n", false, "don't read the configuration file")
	flag.BoolVar(&listSrv, "l", false, "list available servers")
	flag.BoolVar(&html, "H", false, "require random server to support html")
	flag.StringVar(&server, "s", "", "use this server to upload")
	flag.Parse()

	switch {
	case listSrv:
		list()
	default:
		upload(flag.Args())
	}
}
