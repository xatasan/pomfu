package pomfu

import (
	"fmt"
	"io"
	"net/url"
	"os"
)

// A Request holds multiple uploads ready to be uploaded when .Upload is
// called. Each new file has to be declared by calling .Name with it's
// corresponding name. It automatically closes the previous buffer
type Request struct {
	name  []string
	out   []io.Reader
	delay []func(Response) Request
	subrq []struct {
		pomf   *Pomf
		reqest Request
	}
	prev    Response
	minsize int64
}

// The single result of a upload in a UploadInfo struct. This contains
// the destination, the calculated hash-sum and the size in bytes
type UploadInfo struct {
	Url  *url.URL
	Hash string
	Size int
}

// A Response object, that aliases a map, that connects names to upload
// information.
type Response map[string]*UploadInfo

// AddFile is a shorthand for easily adding files to a request. It
// automatically closes the previous buffer by calling .Next, and it
// doesn't allow further data to be appended on via Write
func (r Request) AddFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	r.minsize += stat.Size()
	r.AddReader(file.Name(), file)
	return nil
}

// AddReader opens a new buffer to read in content via .Read. It has to be
// called at least once, before one starts using it.
func (r Request) AddReader(name string, in io.Reader) {
	r.out = append(r.out, in)
	r.name = append(r.name, name)
}

// Upload processed all the input and sends it to a server as specified
// with the method's argument. After finishing, successfully or
// unsuccessfully, it empties it's buffers.
//
// Since a server is randomly chosen, one can specify conditions such as
// the minimum upload size permitted and whether HTML uploads are
// allowed. But if the conditions are too strict, Upload might fail!
func (r Request) Upload(html bool, minsize int64) (Response, error) {
	if minsize <= 0 {
		minsize = r.minsize
	}
	p := getRndServer(html, minsize, 0)
	if p == nil {
		return nil, fmt.Errorf("Failed to choose a random server")
	}
	return p.upload(r)
}

// In case one has to manually choose what server to upload the request
// to, the user can manually pass a reference to Pomf struct, instead of
// letting Pomfu randomly choose a server.
func (r Request) UploadTo(p *Pomf) (Response, error) {
	return p.upload(r)
}

// The Upload function offers a simple method for uploading an io.Reader
// to a random Pomf server.
func Upload(name string, in io.Reader) (Response, error) {
	var req Request
	req.AddReader(name, in)
	return req.Upload(false, 0)
}
