package pomfu

import (
	"testing"

	"bytes"
	"fmt"
	"strings"
)

var tests = []struct {
	input    string
	response map[string]*Pomf
	toerr    bool
}{
	{}, {
		input: `


[server]
maxsize=300000
upload=http://website.com/upload.php
`,
		response: map[string]*Pomf{
			"server": {
				Name:    "server",
				MaxSize: 300000,
				Upload:  mustParseURL("http://website.com/upload.php"),
			},
		},
	}, {
		input: `[server1]
upload=http//server1.org/pomf/upload
maxsize=232000
off=true
name=Server One

[server2]
max=100000000
email=john@servertwo.net
upload=http://serverstwo.net/u
# comment in the middle
 html=true
[server3]
uploadto=http://thirdserver.com/upload.php
about=http://thirdserver.com/about.htm
webmaster=http://thirdserver.com/contact.htm
maximum=30000


`,
		response: map[string]*Pomf{
			"server1": {
				Upload:   mustParseURL("http//server1.org/pomf/upload"),
				MaxSize:  232000,
				Disabled: true,
				Name:     "Server One",
			},
			"server2": {
				Name:        "server2",
				MaxSize:     100000000,
				Email:       mustParseAddress("john@servertwo.net"),
				Upload:      mustParseURL("http://serverstwo.net/u"),
				HtmlAllowed: true,
			},
			"server3": {
				Name:      "server3",
				Upload:    mustParseURL("http://thirdserver.com/upload.php"),
				About:     mustParseURL("http://thirdserver.com/about.htm"),
				Webmaster: mustParseURL("http://thirdserver.com/contact.htm"),
				MaxSize:   30000,
			},
		},
	},
}

func samePomf(p1 *Pomf, p2 *Pomf) bool {
	sEq := func(str1, str2 string) bool { // sanitze
		return strings.ToLower(strings.TrimSpace(str1)) == strings.ToLower(strings.TrimSpace(str2))
	}
	SEq := func(str1, str2 fmt.Stringer) bool {
		return (str1 == str2) ||
			(str1 != nil && str2 != nil && sEq(str1.String(), str2.String()))
	}

	return (p1 == p2) || (sEq(p1.Name, p2.Name) &&
		sEq(p1.Owner, p2.Owner) &&
		(p1.HtmlAllowed == p2.HtmlAllowed) &&
		(p1.Disabled == p2.Disabled) &&
		SEq(p1.About, p2.About) &&
		SEq(p1.Upload, p2.Upload) &&
		SEq(p1.Webmaster, p2.Webmaster) &&
		SEq(p1.Email, p2.Email) &&
		(p1.MaxSize == p2.MaxSize))
}

func TestParser(t *testing.T) {
	for i, test := range tests {
		res, err := parseFile(bytes.NewBufferString(test.input))
		if test.toerr && err == nil {
			t.Fatalf("Expected an error in test %d, but didn't receive one", i)
		} else if !test.toerr && err != nil {
			t.Fatalf("Didn't expect an error in test %d, but received one: %s", i, err)
		}

		if len(res) != len(test.response) {
			t.Errorf("Test %d generated generated different results (length: expected %d, got %d)",
				i, len(test.response), len(res))
		}

		for k, v := range res {
			other, ok := test.response[k]
			if !ok {
				t.Fatalf("Test returned Pomf server with key %v that wasn't expected", k)
				continue
			}

			if !samePomf(v, other) {
				t.Errorf("Pomf instances %#v and %#v with the keyword \"%s\" differ",
					v, test.response[k], k)
			}
		}
	}
}
