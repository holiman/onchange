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
	dialer   *net.Dialer
	callback func(string)
	closeCh  chan any
}

func NewTcpWatcher(subject string, fn func(status string)) (Watch, error) {
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
	}, nil
}

func (t *tcpWatcher) Stop() {
	close(t.closeCh)
}
func (t *tcpWatcher) Start() {
	go t.loop()
}

// loop starts monitoring the given tcp file for changes, and invokes fn
// on changes. The implementation checks the port by default once per five minutes.
// A change is whenever the error-return from a connection attempt changes. This
// means that changes in open->closed or closed->filtered will trigger the callback.
func (t *tcpWatcher) loop() {
	waitTimer := time.NewTimer(5 * time.Minute)
	defer waitTimer.Stop()
	status := t.status()
	slog.Info("TCP initial status", "subject", t.subject, "status", status)
	for {
		select {
		case <-waitTimer.C:
			waitTimer.Reset(5 * time.Minute)
		case <-t.closeCh:
			return
		}
		if newStatus := t.status(); newStatus != status {
			slog.Info("Port change detected", "address", t.subject, "status", newStatus, "previous", status)
			status = newStatus
			t.callback(status)
		}
	}
}

func (t *tcpWatcher) status() string {
	conn, err := t.dialer.Dial("tcp", t.subject)
	if err != nil {
		return err.Error()
	}
	conn.Close()
	return "open"
}
