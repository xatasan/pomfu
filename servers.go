package pomfu

import (
	"net/mail"
	"net/url"
)

// The following map contains all servers that pomfu knows about by
// default. The public variable Servers is aliased to this variable, so
// to hide all the internal details from appearing in godoc
//
// If you own one of these servers, and would wish to have
// your server removed (for whatever reason), please contact
// the developers
//
// based on https://github.com/tsudoko/long-live-pomf
var servers = map[string]*Pomf{
	"subgod": {
		Name:        "sub.god.jp",
		HtmlAllowed: true,
		Upload:      mustParseURL("https://sub.god.jp/upload"),
		About:       mustParseURL("https://sub.god.jp/"),
		Disabled:    false,
		MaxSize:     32 * (1 << 20),
		Owner:       "Xatasan",
		Email:       mustParseAddress("Xatasan <xatasan@firemail.cc>"),
		Webmaster:   mustParseURL("https://sub.god.jp/~xat/"),
	},
}

// Servers is a map of all the Pomf servers Pomfu knows without any
// configuration
var Servers = servers

// forces parsing URLs
func mustParseURL(rawurl string) *url.URL {
	U, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return U
}

// forces parsing mail addresses
func mustParseAddress(address string) *mail.Address {
	A, err := mail.ParseAddress(address)
	if err != nil {
		panic(err)
	}
	return A
}
