package pomfu

import (
	"math/rand"
	"net/mail"
	"net/url"
	"time"
)

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

func RandomServer(html bool) *Pomf {
	return getRndServer(html, 0)
}
