package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/xatasan/pomfu"
)

// helper method to add all the arguments to the current request
func uploadArgs(req pomfu.Request, args []string) {
	for _, arg := range args {
		if err := req.AddFile(arg); err != nil {
			switch {
			case os.IsNotExist(err), os.IsPermission(err):
				fmt.Fprintf(os.Stderr, "Couldn't open file \"%s\"", arg)
				os.Exit(0)
			default:
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		}
	}
}

// upload the files corresponding to the arguments passed in the
// args. If args is empty, use the standard input.
func upload(args []string) {
	var req pomfu.Request
	if len(args) == 0 { // no arguments -> read from standard input
		req.AddReader("-", os.Stdin)
	} else { // arguments given -> use helper to add them to the request
		uploadArgs(req, args)
	}

	resp, err := req.UploadTo(srv) // upload request to specified server (see pomfu.go)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while uploading: %s\n", err.Error())
		os.Exit(1)
	}

	if len(args) == 0 { // if no arguments were given, just print out the url
		fmt.Println(resp["-"].Url)
	} else { // ... otherwise produce a table with local file names and URLs
		out := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		defer out.Flush()
		for name, ui := range resp {
			fmt.Fprintf(out, "%s\t%s\n", name, ui.Url)
		}
	}
}
