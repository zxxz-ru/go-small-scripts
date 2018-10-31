package user

import (
	"regexp"
	"testing"
)

func Test_makeData(t *testing.T) {
	r := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)
	cases := []struct {
		name  string
		login string
	}{
		{"first", "sdfaggfdgdfsg"},
		{"second", "@33s&d"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := makeUuid(c.login)
			if !r.MatchString(res) {
				t.Fatalf("%s: %s is not valid uuid", c.name, res)
			}
		})
	}
}
