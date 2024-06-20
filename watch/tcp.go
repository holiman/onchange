// onchange: react on changes
// Copyright 2024 Martin Holst Swende @holiman
// SPDX-License-Identifier: BSD-3-Clause

package watch

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

type tcpWatcher struct {
	subject  string
	callback func(string)
	closeCh  chan any
	waitTime time.Duration

	dialer *net.Dialer
}

func NewTcpWatcher(subject string, waitTime time.Duration, fn func(status string)) (Watcher, error) {
	if _, err := net.ResolveTCPAddr("tcp", subject); err != nil {
		return nil, fmt.Errorf("address resolution error: %v (input %v)", err, subject)
	}
	dialer := &net.Dialer{
		Timeout: time.Second,
	}
	return &tcpWatcher{
		subject:  subject,
		dialer:   dialer,
		callback: fn,
		closeCh:  make(chan any),
		waitTime: waitTime,
	}, nil
}

func (w *tcpWatcher) Stop() {
	close(w.closeCh)
}
func (w *tcpWatcher) Start() {
	go w.loop()
}

// loop starts monitoring the given tcp file for changes, and invokes fn
// on changes. The implementation checks the port by default once per five minutes.
// A change is whenever the error-return from a connection attempt changes. This
// means that changes in open->closed or closed->filtered will trigger the callback.
func (w *tcpWatcher) loop() {
	waitTimer := time.NewTimer(w.waitTime)
	defer waitTimer.Stop()
	status := w.status()
	slog.Info("TCP initial status", "subject", w.subject, "status", status)
	for {
		select {
		case <-waitTimer.C:
			waitTimer.Reset(w.waitTime)
		case <-w.closeCh:
			return
		}
		if newStatus := w.status(); newStatus != status {
			slog.Info("Port change detected", "address", w.subject, "status", newStatus, "previous", status)
			status = newStatus
			w.callback(status)
		}
	}
}

func (w *tcpWatcher) status() string {
	conn, err := w.dialer.Dial("tcp", w.subject)
	if err != nil {
		return err.Error()
	}
	conn.Close()
	return "open"
}
