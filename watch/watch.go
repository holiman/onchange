// onchange: react on changes
// Copyright 2024 Martin Holst Swende @holiman
// SPDX-License-Identifier: BSD-3-Clause

package watch

type Watcher interface {
	Start()
	Stop()
}
