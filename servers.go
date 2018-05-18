package pomfu

import (
	"net/mail"
	"net/url"
)

// internal server key to fall back onto, if pomfu fails to choose a
// random server
const fallback = "subgod"

// This file contains a list of all the pomf servers pomfu
// knows about and will use to upload files.
//
// If you own one of these servers, and would wish to have
// your server removed (for whatever reason), please contact
// the developers
//
// based on https://github.com/tsudoko/long-live-pomf
var Servers = map[string]*Pomf{
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
