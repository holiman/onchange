// onchange: react on changes
// Copyright 2024 Martin Holst Swende @holiman
// SPDX-License-Identifier: BSD-3-Clause

package watch

import (
	"fmt"
	"testing"
)

// This is a not a real test, TODO rewrite
func TestWeb(t *testing.T) {
	w := webWatcher{}
	w.subject = "https://www.dn.se:443/"
	fmt.Printf("status: %v\n", w.status())
	w.subject = "https://www.dn.se:9292"
	fmt.Printf("status: %v\n", w.status())
	w.subject = "localhost:8822"
	fmt.Printf("status: %v\n", w.status())
	w.subject = "https://www.dn.se/i/dont"
	fmt.Printf("status: %v\n", w.status())
	w.subject = " not even an url ?? "
	fmt.Printf("status: %v\n", w.status())

	/*
		status: 200:0xcdfd24c1e1d514de0ec698eb63d717b1556a79cdb389d91333aecf68e33ae18e
		status: Get "https://www.dn.se:9292": dial tcp 151.101.245.91:9292: i/o timeout
		status: Get "localhost:8822": unsupported protocol scheme "localhost"
		status: 404:0xbf4b4d0c4afcffa75b583053660a59a08ad5f0167f1231d39629d6564e7aea91
		status: Get "%20not%20even%20an%20url%20?? ": unsupported protocol scheme ""	*/
}
