package models
import (
	"github.com/up-finder/silk.web/app/db"
	"github.com/up-finder/silk.web/app/utils/testutils"
	"testing"
	"fmt"
	"encoding/base64"
)

func TestFetchScript(t *testing.T) {
	cases := []struct {
		domain   string
		config   string
		script   string
		expected string
		err      error
	}{
		{
			domain: "localhost",
			config: `{"secret":123, "test": "test"}`,
			script:
			"function init(x,y) {\n" +
			"return x+y}",
			expected:
			"(function() {\n" +
			`var config = {"secret":123, "test": "test"}` + "\n" +
			"function init(x,y) {\n" +
			"return x+y}}())",
			err:nil,
		},
		{
			domain: "localhost",
			config: ``,
			script:
			"function init(x,y) {\n" +
			"return x+y}",
			expected:
			"(function() {\n" +
			`var config = {}` + "\n" +
			"function init(x,y) {\n" +
			"return x+y}}())",
			err:nil,
		},
		{
			domain: "localhost",
			config: `{"secret":123, "test": "test"}`,
			script:"",
			expected:
			"(function() {\n" +
			`var config = {"secret":123, "test": "test"}` + "\n" +
			"}())",
			err:nil,
		},
		{
			domain: "",
			config: ``,
			err:fmt.Errorf(""),
		},
	}

	for _, c := range cases {
		testutils.Setup()
		db.Redis.HSet(CONFIG_HASH_NAME, c.domain, c.config)
		encoded:=base64.StdEncoding.EncodeToString([]byte(c.script))
		db.Redis.LPush(SCRIPT_QUEUE_NAME, encoded)
		got, err := FetchScript(c.domain)
		if (err == nil) && (c.err == nil) {
			if got.Value != c.expected {
				t.Errorf("TestFetchScript: on input %+v\nexpected:\n---\n%s\n---\n\ngot:\n---\n%s", c, c.expected, got.Value)
			}
		}
		if ((c.err == nil) && (err != nil) || (c.err != nil)&&(err == nil)) {
			t.Errorf("TestFetchScript: on input %+v expected err: %v, got:%v", c, c.err, err)
		}
	}
}
