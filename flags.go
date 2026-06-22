package main

import (
	"fmt"
	"strings"
)

// version is set at build time via -ldflags "-X main.version=…".
var version = "dev"

const helpText = `mozeidon-z-messaging — Mozeidon-Z native-messaging host

A browser native-messaging host: it proxies between the Mozeidon-Z browser
extension (stdin/stdout) and the Mozeidon-Z CLI (Unix-socket IPC). It is
normally launched by the browser, not run directly.

Usage:
  mozeidon-z-messaging [--version] [--help]

Flags:
  -v, --version   print version and exit
  -h, --help      print this help and exit`

// handleFlags inspects process args (os.Args[1:]). If the first arg is a
// recognized flag it returns handled=true plus the text to print. Browser
// launches pass a manifest path / extension id as the first arg, so they
// fall through (handled=false) to the proxy.
func handleFlags(args []string) (handled bool, output string) {
	if len(args) == 0 {
		return false, ""
	}
	switch strings.TrimSpace(args[0]) {
	case "--version", "-v", "version":
		return true, fmt.Sprintf("mozeidon-z-messaging %s", version)
	case "--help", "-h", "help":
		return true, helpText
	default:
		return false, ""
	}
}
