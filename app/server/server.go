package server

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/up-finder/silk.web/app"
	"github.com/up-finder/silk.web/app/server/middleware"
	"net/http"
)

type ServerClass struct {
	Port int
}

var Server = NewServer()

func NewServer() *ServerClass {
	port := app.Env.Server.Port
	return &ServerClass{Port: port}
}

// Starts a new http server on the port spectified by app.Env parameters
func Start() {
	Server.Start()
}

func (s *ServerClass) Start() {
	r := NewRouter()
	stack := middleware.NewStack(r)
	stack.Register(middleware.Cors, middleware.Log, middleware.Domain,
		middleware.Ip, middleware.Auth)
	http.Handle("/", stack.Compile())
	log.Infof("Starting server on port %d", s.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
	if err != nil {
		log.Errorf("Error starting server: %v", err)
	}
}
