package models

import (
	"fmt"
	"github.com/up-finder/silk.web/app/db"
	"github.com/up-finder/silk.web/app/json"
	"github.com/up-finder/silk.web/app/utils/testutils"
	"testing"
)

func TestGetEventStatus(t *testing.T) {
	testutils.Setup()
	db.Redis.HSet(UUID_HASH_NAME, "123", "data")
	cases := []struct {
		uuid     string
		expected EventStatus
	}{
		{
			uuid:     "123",
			expected: EventSaved,
		},
		{
			uuid:     "1235",
			expected: EventNotSaved,
		},
		{
			uuid:     "",
			expected: EventNotSaved,
		},
	}
	for _, c := range cases {
		got := GetEventStatus(c.uuid)
		if got != c.expected {
			t.Errorf("Event Model: GetEventStatus: Expected (%d), got (%d)", c.expected, got)
		}
	}
}

func TestSave(t *testing.T) {
	testutils.Setup()
	db.Redis.HSet(UUID_HASH_NAME, "123", "data")
	cases := []struct {
		uuid     string
		expected bool
	}{
		{
			uuid:     "123",
			expected: false,
		},
		{
			uuid:     "1235",
			expected: true,
		},
		{
			uuid:     "",
			expected: false,
		},
		{
			uuid:     "nil",
			expected: false,
		},
	}
	for _, c := range cases {
		data := fmt.Sprintf(`{"uuid": "%s", "payload":{}}`, c.uuid)
		if c.uuid == "nil" {
			data = fmt.Sprintf(`{"payload":{}}`)
		}
		e := NewEvent(json.JSON(data))
		got := e.Save()
		if got != c.expected {
			t.Errorf("Event Model Save: expected (%t), got (%t)", c.expected, got)
		}
	}
}
