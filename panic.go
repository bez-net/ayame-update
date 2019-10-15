package main

import "runtime/debug"

func getStringPanicStack() (str string) {
	if r := recover(); r != nil {
		str = string(debug.Stack())
	}
	return
}
