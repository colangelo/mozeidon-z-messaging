package main

import (
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestSocketPath(t *testing.T) {
	got := socketPath("mozeidon_native_app_1_abcdef12")
	want := "/tmp/mozeidon_native_app_1_abcdef12.sock"
	if got != want {
		t.Errorf("socketPath = %q, want %q", got, want)
	}
}

func TestPidFromSocketName(t *testing.T) {
	cases := []struct {
		path    string
		wantPid int
		wantOK  bool
	}{
		{"/tmp/mozeidon_native_app_81549_d0adee35.sock", 81549, true},
		{"/tmp/mozeidon_native_app_16311_bcb00359.sock", 16311, true},
		{"mozeidon_native_app_42_aaaaaaaa.sock", 42, true},
		{"/tmp/something_else.sock", 0, false},                // not ours
		{"/tmp/mozeidon_native_app_notapid_x.sock", 0, false}, // pid not numeric
		{"/tmp/mozeidon_native_app_.sock", 0, false},          // no pid segment
	}
	for _, c := range cases {
		pid, ok := pidFromSocketName(c.path)
		if pid != c.wantPid || ok != c.wantOK {
			t.Errorf("pidFromSocketName(%q) = (%d,%v), want (%d,%v)",
				c.path, pid, ok, c.wantPid, c.wantOK)
		}
	}
}

func TestStaleSocketsToRemove(t *testing.T) {
	paths := []string{
		"/tmp/mozeidon_native_app_100_aaaaaaaa.sock", // dead → remove
		"/tmp/mozeidon_native_app_200_bbbbbbbb.sock", // alive sibling → keep
		"/tmp/unrelated.sock",                        // not ours → ignore
	}
	isAlive := func(pid int) bool { return pid == 200 }

	got := staleSocketsToRemove(paths, isAlive)
	want := []string{"/tmp/mozeidon_native_app_100_aaaaaaaa.sock"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("staleSocketsToRemove = %v, want %v", got, want)
	}
}

// TestSweepStaleSockets exercises the real glob+remove against /tmp: a socket
// owned by a dead pid is removed, one owned by a live pid (us) is kept.
func TestSweepStaleSockets(t *testing.T) {
	const deadPid = 999999 // above macOS PID_MAX → guaranteed not running
	alivePid := os.Getpid()

	dead := socketPath("mozeidon_native_app_" + strconv.Itoa(deadPid) + "_sweeptst")
	alive := socketPath("mozeidon_native_app_" + strconv.Itoa(alivePid) + "_sweeptst")
	for _, p := range []string{dead, alive} {
		if err := os.WriteFile(p, nil, 0o600); err != nil {
			t.Fatalf("setup %s: %v", p, err)
		}
		defer os.Remove(p)
	}

	sweepStaleSockets()

	if _, err := os.Stat(dead); !os.IsNotExist(err) {
		t.Errorf("dead-pid socket should be removed, stat err = %v", err)
	}
	if _, err := os.Stat(alive); err != nil {
		t.Errorf("live-pid socket should be kept, stat err = %v", err)
	}
}
