package utils

import "time"

// Returns a formatted string of the current time
var Timestamp = func() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05.999 MST")
}
