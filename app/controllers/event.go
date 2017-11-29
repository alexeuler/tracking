// Controllers package contains MVC controllers
package controllers

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/up-finder/silk.web/app/json"
	"github.com/up-finder/silk.web/app/models"
	"github.com/up-finder/silk.web/app/server/middleware"
	utils "github.com/up-finder/silk.web/app/utils"
	"net/http"
)

// Singleton for the event controller functions container
var Event = NewEvent(eventShow, eventCreate)

// Class for the event controller functions container
type EventClass struct {
	Show   http.HandlerFunc
	Create http.HandlerFunc
}

// Event controller constructor - show and create functions are required
func NewEvent(show, create http.HandlerFunc) *EventClass {
	return &EventClass{Show: show, Create: create}
}

// Gets the status of the event by the id parameter, passed in the route
func eventShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string = vars["id"]
	switch models.GetEventStatus(id) {
	case models.EventSaved:
		w.Write([]byte(`{"saved":true}`))
	case models.EventNotSaved:
		w.Write([]byte(`{"saved":false}`))
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Gets the data to save in request body, adds timestamp, ip address, domain and user agent and saves to file db
func eventCreate(w http.ResponseWriter, r *http.Request) {
	body := utils.ReadBody(&r.Body)
	data := json.JSON(body)
	attrs := fmt.Sprintf(`{"created_at":"%s", "ip_address":"%s", "domain":"%s", "fqdn":"%s", "user_agent":"%s"}`,
		utils.Timestamp(), r.Header.Get(middleware.IP_HEADER), r.Header.Get(middleware.DOMAIN_HEADER),
		r.Header.Get(middleware.FQDN_HEADER), r.UserAgent())
	res, err := data.Merge(json.JSON(attrs))
	if err != nil {
		log.Errorf("Merge for Event::Create request with body %v with attrs %v failed: %v", data, attrs, err)
	}
	event := models.NewEvent(res)
	event.Save()
}
