package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/xatasan/pomfu"
)

func list() {
	if !noConf {
		pomfu.Setup()
	}

	out := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(out, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		"key", "name", "about", "max size", "enabled",
		"html", "owner")
	for name, pomf := range pomfu.Servers {
		fmt.Fprintf(out, "%s\t%s\t%s\t%s\t%t\t%t\t%s\n",
			name, pomf.Name, pomf.About, byteSize(pomf.MaxSize), !pomf.Disabled,
			pomf.HtmlAllowed, pomf.Owner.Email)
	}
	out.Flush()
}
