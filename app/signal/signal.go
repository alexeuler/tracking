/*
	This package implements signal handling for the go process
	Signals:
		SIGUSR1 - make export
*/

package signal

import (
	log "github.com/Sirupsen/logrus"
	"github.com/up-finder/silk.web/app/db"
	"os"
	"os/signal"
	"syscall"
)

var channel chan os.Signal

// SIGUSR 1 signal triggers database export
func init() {
	channel = make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGUSR1)
	go func() {
		for {
			sig := <-channel
			log.Debugf("Received signal: %v", sig)
			switch sig {
			case syscall.SIGUSR1:
				sigusr1()
			default:
				log.Errorf("Received unrecognized signal %v", sig)
			}
		}
	}()
}

var sigusr1 = func() {
	db.File.Export()
}
