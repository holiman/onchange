// onchange: react on changes
// Copyright 2024 Martin Holst Swende @holiman
// SPDX-License-Identifier: BSD-3-Clause

package watch

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type webWatcher struct {
	subject  string
	callback func(string)
	closeCh  chan any
}

func NewWebWatcher(subject string, fn func(status string)) (Watch, error) {
	if _, err := url.Parse(subject); err != nil {
		return nil, err
	}
	return &webWatcher{
		subject:  subject,
		callback: fn,
		closeCh:  make(chan any),
	}, nil
}

func (w *webWatcher) Stop() {
	close(w.closeCh)
}
func (w *webWatcher) Start() {
	go w.loop()
}

// loop starts monitoring the HTTP URL, and invokes fn
// on changes. The implementation checks by default once per five minutes.
func (w *webWatcher) loop() {
	waitTimer := time.NewTimer(5 * time.Minute)
	defer waitTimer.Stop()
	status := w.status()
	slog.Info("HTTP initial status", "subject", w.subject, "status", status)
	for {
		select {
		case <-waitTimer.C:
			waitTimer.Reset(5 * time.Minute)
		case <-w.closeCh:
			return
		}
		if newStatus := w.status(); newStatus != status {
			slog.Info("HTTP uri change detected", "address", w.subject,
				"status", newStatus, "previous", status)
			status = newStatus
			w.callback(status)
		}
	}
}

func (w *webWatcher) status() (status string) {
	resp, err := http.Get(w.subject)
	if err != nil {
		status = err.Error()
	} else {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		status = fmt.Sprintf("%d:%#x", resp.StatusCode, sha256.Sum256(data))
	}
	return status
}
