package pomfu

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// A Request holds multiple uploads ready to be uploaded when .Upload is
// called. Each new file has to be declared by calling .Name with it's
// corresponding name. It automatically closes the previous buffer
type Request struct {
	name []string
	out  []io.Reader
}

// AddFile is a shorthand for easily adding files to a request. It
// automatically closes the previous buffer by calling .Next, and it
// doesn't allow further data to be appended on via Write
func (r *Request) AddFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	r.AddReader(file.Name(), file)
	return nil
}

// Next opens a new buffer to read in content via .Read. It has to be
// called at least once, before one starts using it.
func (r *Request) AddReader(name string, in io.Reader) {
	r.out = append(r.out, in)
	r.name = append(r.name, name)
}

// Upload processed all the input and sends it to a server as specified
// with the method's argument. After finishing, successfully or
// unsuccessfully, it empties it's buffers.
//
// If p is null, a random server will be chosen
func (r *Request) Upload(p *Pomf) ([]*Response, error) {
	if r == nil || 0 == len(r.name) {
		return nil, nil
	}

	if p == nil {
		p = RandomServer(false)
	}

	var (
		pr, pw = io.Pipe()
		errch  = make(chan error, 1)
		resch  = make(chan []*Response, 1)
		mimech = make(chan string, 1)
	)

	go func() {
		defer pw.Close()
		mpw := multipart.NewWriter(pw)
		mimech <- mpw.FormDataContentType()

		for i, name := range r.name {
			w, err := mpw.CreateFormFile("files", name)
			if err != nil {
				errch <- err
				return
			}
			_, err = io.Copy(w, r.out[i])
			if err != nil {
				errch <- err
				return
			}
		}

		if err := mpw.Close(); err != nil {
			errch <- err
			return
		}
	}()

	go func() {
		url := p.Upload
		url.Query().Set("output", "json")
		resp, err := http.Post(url.String(), <-mimech, pr)
		if err != nil {
			errch <- err
			return
		}

		dec := json.NewDecoder(resp.Body)
		var data struct {
			Success     bool        `json:"success"`
			Errorcode   int         `json:"errorcode"`
			Description string      `json:"description"`
			Files       []*Response `json:"files"`
		}

		err = dec.Decode(&data)
		if err != nil {
			errch <- err
			return
		}
		if !data.Success {
			errch <- fmt.Errorf("Error while uploading (%d on %s): %s",
				data.Errorcode, p.Name, data.Description)
			return
		}
		for _, f := range data.Files {
			f.Url, err = url.Parse(f.RawUrl)
			// TODO: maybe fix broken or partial URLs?
			if err != nil {
				errch <- err
				return
			}
		}
		resch <- data.Files
	}()

	select {
	case err := <-errch:
		return nil, err
	case res := <-resch:
		return res, nil
	}
}
