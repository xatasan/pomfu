package pomfu

import (
	"net/url"
)

type Response struct {
	Name   string   `json:"name"`
	RawUrl string   `json:"url"`
	Url    *url.URL `json:"-"`
	Hash   string   `json:"hash"`
	Size   int      `json:"size"`
}
