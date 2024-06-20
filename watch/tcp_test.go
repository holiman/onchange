// onchange: react on changes
// Copyright 2024 Martin Holst Swende @holiman
// SPDX-License-Identifier: BSD-3-Clause

package watch

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// This is a not a real test, TODO rewrite
func TestTCP(t *testing.T) {
	w := tcpWatcher{
		dialer: &net.Dialer{
			Timeout: time.Second,
		},
	}
	w.subject = "www.dn.se:80"
	fmt.Printf("status: %v\n", w.status())
	w.subject = "www.dn.se:22"
	fmt.Printf("status: %v\n", w.status())
	w.subject = "localhost:8822"
	fmt.Printf("status: %v\n", w.status())
	w.subject = "www.asdflkweaskfwekskdf:1203"
	fmt.Printf("status: %v\n", w.status())
	w.subject = " not even an url ?? "
	fmt.Printf("status: %v\n", w.status())
}
