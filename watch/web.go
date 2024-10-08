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
	"strings"
	"time"
)

type webWatcher struct {
	subject  string
	content  string
	callback func(string)
	closeCh  chan any
	waitTime time.Duration
}

func NewWebWatcher(subject string, waitTime time.Duration, fn func(status string), contentWatch string) (Watcher, error) {
	if _, err := url.Parse(subject); err != nil {
		return nil, err
	}
	return &webWatcher{
		subject:  subject,
		content:  contentWatch,
		callback: fn,
		closeCh:  make(chan any),
		waitTime: waitTime,
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
	waitTimer := time.NewTimer(w.waitTime)
	defer waitTimer.Stop()
	status := w.status()
	slog.Info("HTTP initial status", "subject", w.subject, "status", status, "polling time", w.waitTime)
	for {
		select {
		case <-waitTimer.C:
			waitTimer.Reset(w.waitTime)
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
	cli := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := cli.Get(w.subject)
	if err != nil {
		return err.Error()
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if w.content == "" {
		return fmt.Sprintf("%d:%#x", resp.StatusCode, sha256.Sum256(data))
	}
	if strings.Contains(string(data), w.content) {
		return fmt.Sprintf("%d:content present", resp.StatusCode)
	}
	return fmt.Sprintf("%d:content missing", resp.StatusCode)
}
