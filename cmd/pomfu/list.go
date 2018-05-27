package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/xatasan/pomfu"
	"math"
)

const unit = 1 << 10

// convert a byte count to a human readable representation
func byteSize(bytes int64) string {
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	b := float64(bytes)
	exp := math.Floor(math.Log(b) / math.Log(unit))
	return fmt.Sprintf("%.2g%ciB",
		b/(math.Pow(unit, exp)),
		"KMGTPE"[int(exp)-1])
}

// print a formatted list of all the servers the pomfu client knows
// about
func list() {
	out := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
	for name, pomf := range pomfu.Servers {
		fmt.Fprintf(out, "%s\t%s\t%s\t%s\t%t\t%t\t%s\n",
			name, pomf.Name, pomf.About, byteSize(pomf.MaxSize), !pomf.Disabled,
			pomf.HtmlAllowed, pomf.Email)
	}
	out.Flush()
}
