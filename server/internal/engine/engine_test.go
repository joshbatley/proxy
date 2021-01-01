package engine

import (
	"fmt"
	"net/url"
	"testing"
)

func TestRemapper(t *testing.T) {

	baseURL := "http://www.test.com/"
	var tests = []struct {
		a, b, c string
		want    string
	}{
		{
			baseURL + "q?s=a_search&p=1",
			baseURL + `q\?s=(.*)\&p=1`,
			baseURL + "a/$1?p=1",
			baseURL + "a/a_search?p=1",
		},
		{
			baseURL + "q?s=a_search&p=10",
			baseURL + `q\?s=(.*)\&p=\d*`,
			baseURL + "a/$1?p=1",
			baseURL + "a/a_search?p=1",
		},
		{
			baseURL + "q?s=a_search&p=13",
			baseURL + `q\?s=(.*)\&p=(\d)`,
			baseURL + "a/$1?p=$2",
			baseURL + "a/a_search?p=13",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s remapped to %s", tt.a, tt.want)
		t.Run(testname, func(t *testing.T) {
			formattedURL, _ := url.ParseRequestURI(tt.a)
			testEngine := RuleEngine{
				url:        formattedURL,
				collection: 1,
				matchedRule: &Rule{
					Pattern:      tt.b,
					Delay:        0,
					Expiry:       0,
					ForceCors:    0,
					SaveResponse: 0,
					SkipOffline:  0,
					RemapRegex:   tt.c,
				},
				rules: make([]Rule, 0),
			}
			ans := testEngine.Remapper()
			if ans.String() != tt.want {
				t.Errorf("got \n %s,\n want \n%s", ans, tt.want)
			}
		})
	}
}
