package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/xatasan/pomfu"
)

// upload the files corresponding to the arguments passed in the
// args. If args is empty, use the standard input.
func upload(args []string) {
	if !noConf {
		pomfu.Setup()
	}

	var (
		req pomfu.Request
		srv *pomfu.Pomf
	)

	if len(args) == 0 {
		req.AddReader("stdin.txt", os.Stdin)
	} else {
		for _, arg := range args {
			err := req.AddFile(arg)
			if err != nil {
				switch {
				case os.IsNotExist(err), os.IsPermission(err):
					fmt.Fprintf(os.Stderr, "Couldn't open file \"%s\"", arg)
					return
				default:
					fmt.Fprintln(os.Stderr, err.Error())
					os.Exit(1)
				}
			}
		}
	}

	if server != "" {
		var ok bool
		srv, ok = pomfu.Servers[server]
		if !ok {
			fmt.Fprintf(os.Stderr, "No server with key: \"%s\"\n", server)
			os.Exit(1)
		}
	}

	resp, err := req.Upload(srv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while uploading: %s\n", err.Error())
		os.Exit(1)
	}

	out := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for name, ui := range resp {
		fmt.Fprintf(out, "%s\t%s\n", name, ui.Url)
	}
	out.Flush()
}
