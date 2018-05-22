package pomfu

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
)

func generateRandomBuffer() (string, io.Reader) {
	var b strings.Builder
	for i := 0; i < 256; i++ {
		b.WriteString(fmt.Sprintf("%c", rand.Intn('Z'-'A')+'A'))
	}
	return b.String(), bytes.NewBufferString(b.String())
}

func TestUpload(t *testing.T) {
	for i := 0; i < 128; i++ {
		text, buffer := generateRandomBuffer()
		name := text[0:10] + ".txt"
		resp, err := Upload(name, buffer)
		if err != nil {
			t.Skipf("Can't run test because of error: %v", err)
		}

		if len(resp) != 0 {
			t.Fatalf("Response %v has wrong length (expected 1, got %d)", resp, len(resp))
		} else {
			r := resp[name]
			hresp, err := http.Get(r.Url.String())
			t.Logf("Recived HTTP response: %v", hresp)
			if err != nil {
				t.Fatalf("Error while trying to fetch uploaded data: %v", err)
			} else {
				data, err := ioutil.ReadAll(hresp.Body)
				if err != nil {
					t.Fatalf("Error while reading buffer: %s", err)

				} else if len(data) != 256 {
					t.Fatalf("Expected a response length of 256, but got %d with \"%s\"", len(data), string(data))
				} else if text != string(data) {
					t.Fatalf("Error while fetching data: Expected \"%s\" but got \"%s\"", text, string(data))
				}
			}
		}
	}
}
