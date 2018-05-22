package main // most of this code is based on the registars source code

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/xatasan/pomfu"
)

var templ = template.New(`
<!DOCTYPE html>
<meta name="viewport" content="width=device-width" />
<meta charset="utf-8" />
<table><thead>
<tr><th>File</th><th>Name</th><th>Hashcode</th><th>Size</th></tr>
</thead><tbody>
{{ range $name, $ui := . }}
<tr>
	<td><a href="{{ $ui.Url }}">{{ $ui.Url }}</a></td>
	<td><code>{{ $name }}</code></td>
	<td><code>{{ printf "%.10s" $ui.Hash }}...</code></td>
	<td>{{ $ui.Size }} bytes</td>
</tr>
{{ end }}
</tbody></table>`)

func upload(w http.ResponseWriter, r *http.Request) {
	html := r.URL.Query().Get("html") != ""
	mpr, err := r.MultipartReader()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	form, err := mpr.ReadForm(1 << 32)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	//// START MAIN EXAMPLE CODE ////
	if !noConf { // prevent pomfu from reading it's config if wished
		pomfu.ReadConfig()
	}

	var request pomfu.Request // a request object is created
	for _, file := range form.File["files"] {
		fh, err := file.Open()
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		request.AddReader(file.Filename, fh) // each new form file is added as a io.Reader
	}

	resp, err := request.Upload(html, limit) // the request is submitted
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	//// END MAIN EXAMPLE CODE   ////

	// NOTE: to keep it brief, not all "required" output types have been
	// implemented here, such as csv, gyazo and text.
	//
	// This is a feature, not a bug
	switch r.URL.Query().Get("output") {
	case "html":
		w.Header().Set("Content-Type", "text/html")
		templ.ExecuteTemplate(w, "", resp)
	default: // "json"
		w.Header().Set("Content-Type", "application/json")
		type File struct {
			Name string `json:"name"`
			Url  string `json:"url"`
			Hash string `json:"hash"`
			Size int    `json:"size"`
		}

		var files []File
		for name, ui := range resp {
			files = append(files, File{name, ui.Url.String(), ui.Hash, ui.Size})
		}

		json.NewEncoder(w).Encode(struct {
			Success     bool   `json:"success"`
			Errorcode   int    `json:"errorcode"`
			Description string `json:"description"`
			Files       []File `json:"files"`
		}{true, 200, "", files})
	}
}
