package main

import (
	"strings"
	"testing"
)

func TestHandleFlags_Version(t *testing.T) {
	version = "9.9.9"
	handled, out := handleFlags([]string{"--version"})
	if !handled {
		t.Fatal("expected --version to be handled")
	}
	if out != "mozeidon-z-messaging 9.9.9" {
		t.Fatalf("got %q", out)
	}
}

func TestHandleFlags_VersionShortAlias(t *testing.T) {
	if handled, _ := handleFlags([]string{"-v"}); !handled {
		t.Fatal("expected -v to be handled")
	}
}

func TestHandleFlags_Help(t *testing.T) {
	handled, out := handleFlags([]string{"-h"})
	if !handled {
		t.Fatal("expected -h to be handled")
	}
	if !strings.Contains(out, "native-messaging host") {
		t.Fatalf("help text missing expected phrase: %q", out)
	}
}

func TestHandleFlags_BrowserArgFallsThrough(t *testing.T) {
	// Firefox launches the host with the manifest path + extension id.
	handled, _ := handleFlags([]string{"/path/to/mozeidon.json", "mozeidon-z@a-layer.io"})
	if handled {
		t.Fatal("browser launch args must fall through to the proxy")
	}
}

func TestHandleFlags_NoArgs(t *testing.T) {
	if handled, _ := handleFlags(nil); handled {
		t.Fatal("no args must fall through to the proxy")
	}
}
