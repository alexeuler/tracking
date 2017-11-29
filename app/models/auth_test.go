package models

import (
	"github.com/up-finder/silk.web/app/db"
	"github.com/up-finder/silk.web/app/utils/testutils"
	"testing"
)

func TestFetchAuth(t *testing.T) {
	testutils.Setup()
	db.Redis.HSet(AUTH_HASH_NAME, "domain", "test")
	auth := FetchAuth("domain")
	if auth.Key != "test" {
		t.Errorf("Test Auth Model: expected (test) value, got: %s", auth.Key)
	}
}
