package pomfu

import (
	"net/mail"
	"net/url"
)

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
		"sub.god.jp",
		true,
		mustParseURL("https://sub.god.jp/upload"),
		mustParseURL("https://sub.god.jp/"),
		false,
		32 * (1 << 20),
		Owner{
			"Xatasan",
			mustParseAddress("Xatasan <xatasan@firemail.cc>"),
			mustParseURL("https://sub.god.jp/~xat/"),
		},
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
