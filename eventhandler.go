// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package eventsourcing

type EventHandler interface {
	Handle(EventMessage)
}
