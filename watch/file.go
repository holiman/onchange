// onchange: react on changes
// Copyright 2024 Martin Holst Swende @holiman
// SPDX-License-Identifier: BSD-3-Clause

package watch

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type fileWatcher struct {
	subject  string
	callback func(string)
	w        *fsnotify.Watcher
}

func NewFileWatcher(subject string, fn func(string)) (Watcher, error) {
	if st, err := os.Lstat(subject); err != nil {
		return nil, err
	} else if st.IsDir() {
		return nil, fmt.Errorf("%q is a directory, not a file", subject)
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Failed to start filesystem watcher", "err", err)
		return nil, err
	}
	// Monitor the directory, not the file itself.
	if err = w.Add(filepath.Dir(subject)); err != nil {
		return nil, fmt.Errorf("%q: %s", subject, err)
	}
	return &fileWatcher{
		subject:  subject,
		callback: fn,
		w:        w,
	}, nil
}

func (f *fileWatcher) Stop() {
	f.w.Close()
}
func (f *fileWatcher) Start() {
	go f.loop()
}

// File starts monitoring the given file for changes, and invokes fn
// on changes, at a max rate of once per 500 ms.
func (f *fileWatcher) loop() {
	// Wait for file system events and reload.
	// When an event occurs, the reload call is delayed a bit so that
	// multiple events arriving quickly only cause a single reload.
	var (
		debounceDuration = 500 * time.Millisecond
		rescanTriggered  = false
		debounce         = time.NewTimer(0)
	)
	// Ignore initial trigger
	if !debounce.Stop() {
		<-debounce.C
	}
	defer debounce.Stop()
	for {
		select {
		case err, ok := <-f.w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			slog.Error("Error in watcher", "error", err)
		// Read from Events.
		case e, ok := <-f.w.Events:
			if !ok {
				return
			}
			if f.subject == e.Name {
				// Trigger the scan (with delay), if not already triggered
				if !rescanTriggered {
					debounce.Reset(debounceDuration)
					rescanTriggered = true
				}
			}
		case <-debounce.C:
			f.callback("")
			rescanTriggered = false
		}
	}
}
