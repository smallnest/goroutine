// Copyright ©2020 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gid provides functions that will return the runtime's
// ID number for the calling goroutine or its creator.
//
// The implementation is derived from Laevus Dexter's comment in Gophers' Slack #darkarts,
// https://gophers.slack.com/archives/C1C1YSQBT/p1593885226448300 post which linked to
// this playground snippet https://play.golang.org/p/CSOp9wyzydP.
package gid

import (
	"reflect"
	"unsafe"
)

// ID returns the runtime ID of the calling goroutine.
func ID() int64 {
	return idOf(getg(), goidoff)
}

func idOf(g unsafe.Pointer, off uintptr) int64 {
	return *(*int64)(add(g, off))
}

//go:nosplit
func getg() unsafe.Pointer {
	return *(*unsafe.Pointer)(add(getm(), curgoff))
}

//go:linkname runtime_procPin runtime.procPin
func runtime_procPin() int

//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin() int

// PID returns the "P" ID of the calling goroutine.
func PID() int {
	pid := runtime_procPin()
	runtime_procUnpin()
	return pid
}

//go:linkname add runtime.add
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer

//go:linkname getm runtime.getm
//go:nosplit
func getm() unsafe.Pointer

var (
	curgoff = offset("*runtime.m", "curg")
	goidoff = offset("*runtime.g", "goid")
)

// offset returns the offset into typ for the given field.
func offset(typ, field string) uintptr {
	rt := toType(typesByString(typ)[0])
	f, _ := rt.Elem().FieldByName(field)
	return f.Offset
}

//go:linkname typesByString reflect.typesByString
func typesByString(s string) []unsafe.Pointer

//go:linkname toType reflect.toType
func toType(t unsafe.Pointer) reflect.Type
