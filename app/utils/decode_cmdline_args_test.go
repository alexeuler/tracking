package utils

import "testing"

func TestDecodeCmdLineArgs(t *testing.T) {

	equalHashes := func(expected, got map[string]string) bool {
		if len(expected) == len(got) {
			for k, v := range expected {
				if g, ok := got[k]; ok {
					if g == v {
						continue
					}
				}
				return false
			}
			return true
		}
		return false
	}

	testCases := []struct {
		Args     []string
		Expected map[string]string
	}{
		{nil, map[string]string{}},
		{[]string{}, map[string]string{}},
		{[]string{"-v"}, map[string]string{"v": ""}},
		{[]string{"-e", "production"}, map[string]string{"e": "production"}},
		{[]string{"pr", "reg", "erg"}, map[string]string{}},
		{[]string{"-v", "-e", "production"}, map[string]string{"v": "", "e": "production"}},
		{[]string{"silk.web", "-v", "true", "false", "-e", "production"}, map[string]string{"v": "true", "e": "production"}},
	}

	for _, test := range testCases {
		got := DecodeCmdLineArgs(test.Args)
		if !equalHashes(got, test.Expected) {
			t.Errorf("On input: %v expected: %v, got %v", test.Args, test.Expected, got)
		}
	}
}
