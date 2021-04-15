// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// package eventsourcing provides a CQRS reference implementation.
//
// The implementation follows as much as possible the classic reference implementation
// m-r by Greg Young.
//
// The implmentation differs in a number of respects becasue the original is written
// in C# and uses Generics where generics are not available in Go.
// This implementation instead uses interfaces to deal with types in a generic manner
// and used delegate functions to instantiate specific types.
package eventsourcing

import (
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"
)

// typeOf is a convenience function that returns the name of a type
//
// This is used so commonly throughout the code that it is better to
// have this convenience function and also allows for changing the scheme
// used for the type name more easily if desired.
func typeOf(i interface{}) string {
	t := reflect.TypeOf(i)
	if t == nil {
		panic(fmt.Sprintf("Type not found for: %T\n", i))
	}

	return t.Elem().Name()
}

// NewUUID returns a new v4 uuid as a string
func NewUUID() string {
	uuid, error := uuid.NewV4()
	if error != nil {
		panic("Could not generate UUID")
	}
	return uuid.String()
}

// Int64 returns a pointer to int64.
//
// There are a number of places where a pointer to int64
// is required such as expectedVersion argument on the repository
// and this helper function makes keeps the code cleaner in these
// cases.
func Int64(i int64) *int64 {
	return &i
}
