package models
import (
	"github.com/up-finder/silk.web/app/db"
	"bytes"
	"fmt"
	"encoding/base64"
)

const (
	SCRIPT_QUEUE_NAME = "silk:js:script" //the basic script stored in Redis
	CONFIG_HASH_NAME = "silk:js:config" //the per-client config variables json
)

type Script struct {
	Domain string
	Value  string //final js script
}

// Get script for specified domain
// Returns err if domain is empty
// If config is empty it defaults to {}
var FetchScript = func(domain string) (*Script, error) {
	if (domain == "") {
		return nil, fmt.Errorf("Fetch Script: Empty domain")
	}
	conn := db.Redis
	var buf bytes.Buffer
	buf.WriteString("(function() {\n")
	buf.WriteString("var config = ")
	conf, err := conn.HGet(CONFIG_HASH_NAME, domain).Result()
	if err != nil {
		return nil, fmt.Errorf("Fetch Script: %s: %v", CONFIG_HASH_NAME, err)
	}
	if (conf == "") {
		conf = "{}"
	}
	buf.WriteString(conf)
	buf.WriteString("\n")
	scr, err := conn.LIndex(SCRIPT_QUEUE_NAME, 0).Result()
	if err != nil {
		return nil, fmt.Errorf("Fetch Script: %s: %v", SCRIPT_QUEUE_NAME, err)
	}
	decoded, err := base64.StdEncoding.DecodeString(scr)
	if err != nil {
		return nil, fmt.Errorf("Fetch Script: Base64 decoding error: %s: %v", SCRIPT_QUEUE_NAME, err)
	}
	buf.Write(decoded)
	buf.WriteString("}())")
	return &Script{Domain: domain, Value: buf.String()}, nil
}
