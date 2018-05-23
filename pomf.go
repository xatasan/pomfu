package pomfu

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/url"
	"time"
)

// Pomf corresponds to a server that implements the Pomf-API. It
// maintains all the metadata and related information that a user might
// want or need to know about a server.
//
// The fields should be coincided constants, and not changed.
type Pomf struct {
	Name        string
	HtmlAllowed bool
	Upload      *url.URL
	About       *url.URL
	Disabled    bool
	MaxSize     int64
	Owner       string
	Email       *mail.Address
	Webmaster   *url.URL
}

// internal method to choose a random server from the Servers list
func getRndServer(html bool, minsize int64, tries uint) *Pomf {
	if tries >= 20 || len(servers) == 0 {
		return nil
	}

	var (
		rsrc   = rand.New(rand.NewSource(time.Now().UnixNano()))
		chance = 1 / float64(len(servers))
	)

	for _, v := range servers {
		if (html && !v.HtmlAllowed) || v.Disabled || (minsize > 0 && (v.MaxSize < minsize)) {
			chance = 1 / ((1 / chance) - 1)
			continue
		}
		c := rsrc.Float64()
		if c < chance {
			return v
		}
	}

	return getRndServer(html, minsize, tries+1)
}

// internal method to actually upload a request
func (p *Pomf) upload(r Request) (Response, error) {
	var (
		pr, pw = io.Pipe()
		errch  = make(chan error, 1)
		resch  = make(chan Response, 1)
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
		dres, err := http.Post(url.String(), <-mimech, pr)
		if err != nil {
			errch <- err
			return
		}

		dec := json.NewDecoder(dres.Body)
		var data struct {
			Success     bool   `json:"success"`
			Errorcode   int    `json:"errorcode"`
			Description string `json:"description"`
			Files       []struct {
				Name   string `json:"name"`
				RawUrl string `json:"url"`
				Hash   string `json:"hash"`
				Size   int    `json:"size"`
			} `json:"files"`
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

		response := make(Response)

		for _, f := range data.Files {
			url, err = url.Parse(f.RawUrl)
			// TODO: maybe fix broken or partial URLs?
			if err != nil {
				errch <- err
				return
			}

			response[f.Name] = &UploadInfo{
				Size: f.Size,
				Hash: f.Hash,
				Url:  url,
			}
		}
		resch <- response
	}()

	select {
	case err := <-errch:
		return nil, err
	case res := <-resch:
		r.prev = res
		dres, err := r.processDelays(p)
		if err != nil {
			return nil, err
		}
		sreq, err := r.processSubreq(merge(res, dres))
		if err != nil {
			return nil, err
		}

		for _, r := range sreq {
			res = merge(res, r)
		}

		return res, nil
	}
}
