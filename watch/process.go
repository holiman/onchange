package watch

import (
	"log/slog"
	"os"
	"syscall"
	"time"
)

type procWatcher struct {
	pid      int
	callback func(string)
	closeCh  chan any
	waitTime time.Duration
}

func NewProcWatcher(pid int, waitTime time.Duration, fn func(status string)) (Watcher, error) {
	return &procWatcher{
		pid:      pid,
		callback: fn,
		closeCh:  make(chan any),
		waitTime: waitTime,
	}, nil
}

func (w *procWatcher) Stop() {
	close(w.closeCh)
}
func (w *procWatcher) Start() {
	go w.loop()
}

// loop starts monitoring the HTTP URL, and invokes fn
// on changes. The implementation checks by default once per five minutes.
func (w *procWatcher) loop() {
	waitTimer := time.NewTimer(w.waitTime)
	defer waitTimer.Stop()
	status := w.status()
	slog.Info("Proc initial status", "pid", w.pid, "status", status, "polling time", w.waitTime)
	for {
		select {
		case <-waitTimer.C:
			waitTimer.Reset(w.waitTime)
		case <-w.closeCh:
			return
		}
		if newStatus := w.status(); newStatus != status {
			slog.Info("Proc change detected", "pid", w.pid,
				"status", newStatus, "previous", status)
			status = newStatus
			w.callback(status)
		}
	}
}

func (w *procWatcher) status() (status string) {
	proc, err := os.FindProcess(w.pid)
	if err != nil {
		return err.Error()
	}
	// From docs:
	// On Unix systems, FindProcess always succeeds and returns a Process
	// for the given pid, regardless of whether the process exists. To test whether
	// the process actually exists, see whether p.Signal(syscall.Signal(0)) reports
	// an error.
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return err.Error()
	}
	return "running"
}
