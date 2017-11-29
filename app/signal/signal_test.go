package signal

import (
	"syscall"
	"testing"
	"time"
)

//syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
func TestSignal(t *testing.T) {
	ssigusr1 := sigusr1
	defer func() { sigusr1 = ssigusr1 }()
	success := false
	sigusr1 = func() { success = true }
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(time.Millisecond * 10)
	if !success {
		t.Fatalf("SIGUSR1 callback was not launched")
	}
}
