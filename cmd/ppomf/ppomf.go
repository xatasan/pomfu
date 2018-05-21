package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/xatasan/pomfu"
)

const index = `<!DOCTYPE html>
<title>A Pomf Proxy</title>
<meta name="viewport" content="width=device-width" />
<meta charset="utf-8" />
<style>
body{ font-family: arial, sans-serif; max-width: 40em; margin: 1em auto; }
</style>
<body>
<h1>A Pomf Proxy: <em>ppomf</em></h1>
<form action="/upload?output=html" method="post" enctype="multipart/form-data">
<input type="file" name="files" reqired multiple />
<input type="submit" value="Upload" style="float: right;" />
<br/>
HTML: <input type="checkbox" name="html" />
</form>
<p>
This server is a <em>pomf proxy</em>. That means that every file you upload
to this server will be randomly redirected to a <q>real</q> pomf server,
and it's results will be sent back to the <q>real</q> client.
<p>
<strong>Note:</strong> This is a proof-of-concept and a demonstration on how
to use the Go <a href="https://sub.god.jp/~xat/pomfu">Pomfu</a> library. This
shoudln't <em>actually</em> be used in practice!
`

func main() {
	var noConf bool
	var addr string

	flag.BoolVar(&noConf, "n", false, "prevent ppomf from reading the pomfu configuration")
	flag.StringVar(&addr, "a", ":8080", "address to listen on")
	flag.Parse()

	//// START MAIN EXAMPLE CODE ////
	if !noConf {
		pomfu.Setup()
	}
	//// END MAIN EXAMPLE CODE   ////

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/upload.php", upload)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, index)
	})
	fmt.Println("Listening on", addr)
	http.ListenAndServe(addr, nil)
}
