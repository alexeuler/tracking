package controllers

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/up-finder/silk.web/app/json"
	"github.com/up-finder/silk.web/app/models"
	"github.com/up-finder/silk.web/app/server/middleware"
	"github.com/up-finder/silk.web/app/utils"
	"github.com/up-finder/silk.web/app/utils/testutils"
	"net/http"
	"strconv"
	"testing"
)

type contextKey int

const (
	varsKey contextKey = iota
	routeKey
)

func TestEventShow(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/events/{id}").Handler(http.HandlerFunc(eventShow))

	cases := []struct {
		id       string
		expected string
	}{
		{
			id:       "1", // <=100 - not saved, >=100 - saved
			expected: `{"saved":false}`,
		},
		{
			id:       "105",
			expected: `{"saved":true}`,
		},
	}

	for _, c := range cases {
		id := c.id
		r, _ := http.NewRequest("GET", "http://localhost:3000/events/"+id, nil)
		w := testutils.NewFakeResponse()
		saved := models.GetEventStatus
		models.GetEventStatus = func(ide string) models.EventStatus {
			if ide != id {
				t.Errorf("Event Controller: Show: expected id %s, got %s", id, ide)
			}
			i, _ := strconv.Atoi(ide)
			if i > 100 {
				return models.EventSaved
			}
			return models.EventNotSaved
		}
		router.ServeHTTP(w, r)
		models.GetEventStatus = saved
		if w.Body != c.expected {
			t.Errorf("Event Controller: Show: expected response %s, got %s", c.expected, w.Body)
		}
	}
}

var EventStubBuffer bytes.Buffer

type EventStub struct {
}

func NewEventStub(data json.JSON) models.Model {
	EventStubBuffer.Reset()
	EventStubBuffer.Write([]byte(data))
	return &EventStub{}
}

func (e *EventStub) Save() bool {
	return true
}

func TestEventCreate(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path("/events").Handler(http.HandlerFunc(eventCreate))

	type Input struct {
		data       string
		xrealip    string
		xdomain    string
		xfqdn    string
		xuseragent string
		xtime      string
	}
	cases := []struct {
		input    Input
		expected string
	}{
		{
			input: Input{
				data:       `{"uuid":"123", "payload":{"data":1}}`,
				xrealip:    "127.0.0.1",
				xdomain:    "xyz.ru",
				xfqdn:    "xyz.ru",
				xuseragent: `Smith`,
				xtime:      "2016-01-01",
			},
			expected: `{"created_at":"2016-01-01", "ip_address":"127.0.0.1", ` +
				`"domain":"xyz.ru", "fqdn":"xyz.ru", "user_agent":"Smith", "uuid":"123", "payload":{"data":1}}`,
		},
		{
			input: Input{
				data:       `{}`,
				xrealip:    "",
				xdomain:    "",
				xfqdn:    "",
				xuseragent: "",
				xtime:      "",
			},
			expected: `{"created_at":"", "ip_address":"", ` +
				`"domain":"", "fqdn":"", "user_agent":""}`,
		},
		{
			input: Input{
				data:       `{}`,
				xrealip:    "nil",
				xdomain:    "nil",
				xfqdn:    "nil",
				xuseragent: "nil",
				xtime:      "",
			},
			expected: `{"created_at":"", "ip_address":"", ` +
				`"domain":"", "fqdn":"", "user_agent":""}`,
		},
	}
	for _, c := range cases {
		savedTs := utils.Timestamp
		utils.Timestamp = func() string { return c.input.xtime }
		savedNE := models.NewEvent
		models.NewEvent = NewEventStub
		body := bytes.NewReader([]byte(c.input.data))
		r, _ := http.NewRequest("POST", "http://localhost/events", body)
		if c.input.xrealip != "nil" {
			r.Header.Set(middleware.IP_HEADER, c.input.xrealip)
		}
		if c.input.xdomain != "nil" {
			r.Header.Set(middleware.DOMAIN_HEADER, c.input.xdomain)
		}
		if c.input.xfqdn != "nil" {
			r.Header.Set(middleware.FQDN_HEADER, c.input.xfqdn)
		}
		if c.input.xuseragent != "nil" {
			r.Header.Set("User-Agent", c.input.xuseragent)
		}
		w := testutils.NewFakeResponse()
		router.ServeHTTP(w, r)
		got := string(EventStubBuffer.Bytes())
		if got != c.expected {
			t.Errorf("Controller Event: case (%+v), expected (%s), got (%s)", c.input, c.expected, got)
		}
		utils.Timestamp = savedTs
		models.NewEvent = savedNE
	}
}
