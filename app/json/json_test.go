package json

import "testing"

func TestMerge(t *testing.T) {
	cases := []struct {
		mergeTo   JSON
		mergeFrom JSON
		expected  JSON
	}{
		{
			mergeTo:   "{}",
			mergeFrom: `{"a": "1"}`,
			expected:  `{"a": "1"}`,
		},
		{
			mergeTo:   `{b:1}`,
			mergeFrom: `{"a":"1"}`,
			expected:  `{"a":"1", b:1}`,
		},
		{
			mergeTo:   `{b:1}`,
			mergeFrom: `{"a":"1", "c":"2"}`,
			expected:  `{"a":"1", "c":"2", b:1}`,
		},
		{
			mergeTo:   `{"b":{"d":2}}`,
			mergeFrom: `{"a":"1", "c":"2", "e":{"x":"y"}}`,
			expected:  `{"a":"1", "c":"2", "e":{"x":"y"}, "b":{"d":2}}`,
		},
	}

	for _, c := range cases {
		got, err := c.mergeTo.Merge(c.mergeFrom)
		if got != c.expected {
			t.Errorf("Test Merge JSON: on input json: %s with attrs: %s expected: %s, got %s ",
				c.mergeTo, c.mergeFrom, c.expected, got)
		}
		if err != nil {
			t.Errorf("Test Merge JSON: on input json: %s with attrs: %s expected no error, but received: %s ",
				c.mergeTo, c.mergeFrom, err)
		}
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		json     JSON
		key      string
		expected string
		ok       bool
	}{
		{
			json:     `{a:123, b:2}`,
			key:      "a",
			expected: "123",
			ok:       true,
		},
		{
			json:     `{a:"123", b:2}`,
			key:      "a",
			expected: "123",
			ok:       true,
		},
		{
			json:     `{a:"", b:2}`,
			key:      "a",
			expected: "",
			ok:       true,
		},
		{
			json:     `{av:"123", cd: {av:2, ff: 3} }`,
			key:      "a",
			expected: "",
			ok:       false,
		},
	}

	for _, c := range cases {
		res, success := c.json.Get(c.key)
		if success != c.ok {
			t.Errorf("Json: Test Get: on json %v, for key %v - expected success to be %v, got %v", c.json, c.key, c.ok, success)
		}
		if c.expected != res {
			t.Errorf("Json: Test Get: on json %v, for key %v - expected value to be %v, got %v", c.json, c.key, c.expected, res)
		}
	}
}
