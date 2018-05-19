package pomfu

import (
	"math/rand"
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
	MaxSize     int
	Owner       string
	Email       *mail.Address
	Webmaster   *url.URL
}

// internal method to choose a random server from the Servers list
func getRndServer(html bool, tries int) *Pomf {
	if tries >= 20 || len(Servers) == 0 {
		return Servers[fallback] // default server
	}

	var (
		rsrc   = rand.New(rand.NewSource(time.Now().UnixNano()))
		chance = 1 / float64(len(Servers))
	)

	for _, v := range Servers {
		if (html && !v.HtmlAllowed) || v.Disabled {
			chance = 1 / ((1 / chance) - 1)
			continue
		}
		c := rsrc.Float64()
		if c < chance {
			return v
		}
	}

	return getRndServer(html, tries+1)
}

// RandomServer returns a random Server from the Servers list,
// optionally limited to only those servers that support HTML uploads
func RandomServer(html bool) *Pomf {
	return getRndServer(html, 0)
}
