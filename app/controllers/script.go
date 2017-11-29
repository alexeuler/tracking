package controllers
import (
	"net/http"
	"github.com/up-finder/silk.web/app/models"
	"github.com/up-finder/silk.web/app/server/middleware"
	log "github.com/Sirupsen/logrus"
)

// Singleton for the script controller functions container
var Script = NewScript(scriptShow)

// Class for the script controller functions container
type ScriptClass struct {
	Show   http.HandlerFunc
}

// Script controller constructor - show function is required
func NewScript(show http.HandlerFunc) *ScriptClass {
	return &ScriptClass{Show: show}
}

// Gets the script for the domain specified in the domain header
func scriptShow(w http.ResponseWriter, r *http.Request) {
	domain:=r.Header.Get(middleware.DOMAIN_HEADER)
	script, err:=models.FetchScript(domain)
	if err!=nil {
		log.Errorf("Script Show: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(script.Value))
}
