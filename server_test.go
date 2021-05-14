// +build acceptance

package main_test

import (
	"flag"
	"testing"

	"github.com/go-rod/rod"
)

var address = flag.String("address", "localhost:18080", "Address to the server to test")

func TestUI(t *testing.T) {
	page := rod.New().MustConnect().MustPage("http://" + *address).MustWaitLoad()
	attribute := page.MustElement("#app > img:nth-child(1)").MustAttribute("alt")
	if *attribute != "Vue logo" {
		t.Errorf("Attribute 'alt' with wrong value '%s'", *attribute)
	}
}
