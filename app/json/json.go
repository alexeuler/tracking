//Package for working with json as string
package json

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"regexp"
)

type JSON string

// Gets the value by key from JSON string
// Warning! - this method fails if there are two keys in the json, albeit on different depth
// Warning! - doesn't allow for quotes usage inside value
func (j JSON) Get(key string) (string, bool) {
	pattern := fmt.Sprintf(`\"?%s\"?:\s*\"?([^,}\"]*)[,}\"]`, key)
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(j), -1)
	if len(matches) > 1 {
		log.Errorf("JSON: Trying to get the key %v of %v, but found multiple occurences", key, j)
		return "", false
	}
	if len(matches) == 0 {
		return "", false
	}
	res := matches[0][1]
	return res, true
}

// Merges data into JSON and returns new merged JSON
func (j JSON) Merge(data JSON) (res JSON, err error) {
	ok, err := j.Test()
	if !ok {
		return "", fmt.Errorf("Json: Merge: %v", err)
	}
	ok, err = data.Test()
	if !ok {
		return "", fmt.Errorf("Json: Merge: %v", err)
	}
	var buf bytes.Buffer
	for i, r := range j {
		buf.WriteRune(r)
		if i == 0 && r == '{' {
			for j, r2 := range data {
				if (j != 0) && (j != len(data)-1) {
					buf.WriteRune(r2)
				}
			}
			if len(j) != 2 { // nonempty json, i.e. != {}
				buf.WriteString(", ") // add ", "
			}
		}
	}
	return JSON(buf.String()), nil
}

// Tests for correctness of JSON string
func (j JSON) Test() (res bool, err error) {
	str := string(j)
	matched, err := regexp.MatchString(`^\{(.*)\}$`, str)
	if !matched {
		return false, fmt.Errorf("expected {...} string as argument, but received %v", j)
	}
	if err != nil {
		return false, fmt.Errorf("unexpected regexp error testing object: %v", j)
	}
	return true, nil
}
