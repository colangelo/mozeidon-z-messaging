package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// socketDir / socketSuffix mirror golang-ipc v1.2.4 (connect_other.go), which
// hardcodes the unix socket at /tmp/<ipcName>.sock. We replicate the path here so
// we can clean up after ourselves: the library's Server.Close unlinks the socket,
// but only if Close is actually called — on a hard kill it never is, leaking the
// file. (If golang-ipc ever changes its socket location, update this.)
const (
	socketDir     = "/tmp"
	socketSuffix  = ".sock"
	ipcNamePrefix = "mozeidon_native_app_"
)

// socketPath returns the unix-socket path golang-ipc creates for an ipcName.
func socketPath(ipcName string) string {
	return filepath.Join(socketDir, ipcName+socketSuffix)
}

// pidFromSocketName extracts the pid embedded in a mozeidon ipc socket path.
// Names look like: mozeidon_native_app_<pid>_<profile8>.sock
func pidFromSocketName(path string) (int, bool) {
	base := strings.TrimSuffix(filepath.Base(path), socketSuffix)
	if !strings.HasPrefix(base, ipcNamePrefix) {
		return 0, false
	}
	pidStr, _, ok := strings.Cut(strings.TrimPrefix(base, ipcNamePrefix), "_")
	if !ok {
		return 0, false
	}
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		return 0, false
	}
	return pid, true
}

// staleSocketsToRemove selects socket paths whose owning pid is no longer alive.
// Sockets of live instances (e.g. a sibling native-app for another browser
// profile) are kept.
func staleSocketsToRemove(paths []string, isAlive func(int) bool) []string {
	var stale []string
	for _, p := range paths {
		pid, ok := pidFromSocketName(p)
		if !ok {
			continue
		}
		if !isAlive(pid) {
			stale = append(stale, p)
		}
	}
	return stale
}

// pidAlive reports whether a process with the given pid exists (single-user box,
// so EPERM isn't a concern). Signal 0 performs an existence check without
// actually signalling the process.
func pidAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// sweepStaleSockets removes leftover mozeidon ipc sockets from instances that were
// hard-killed (SIGKILL/crash) and so never ran their own shutdown cleanup. Called
// at startup; complements the per-instance cleanup done on graceful shutdown.
func sweepStaleSockets() {
	matches, err := filepath.Glob(filepath.Join(socketDir, ipcNamePrefix+"*"+socketSuffix))
	if err != nil {
		return
	}
	for _, p := range staleSocketsToRemove(matches, pidAlive) {
		os.Remove(p)
	}
}
