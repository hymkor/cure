//go:build !windows
// +build !windows

package main

type _Handle = int

const (
	_ConIn = "/dev/tty"
)
