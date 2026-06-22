# CLAUDE.md

Guidance for Claude Code working in this repo.

## What this is

`mozeidon-z-messaging` — the native-messaging host (browser ⇄ CLI IPC bridge) for the Mozeidon-Z
stack. ~230 lines of Go. Hard fork of `egovelox/mozeidon-native-app` (remote `upstream`).

## Identifiers

- IPC **socket** base name `mozeidon_native_app` (generated form
  `mozeidon_native_app_<pid>_<profileId8>`) — **FROZEN**: it's the contract with the `mozeidon-z` CLI.
- Native-messaging **host name** `"mozeidon_z"` — the extension calls `connectNative("mozeidon_z")`
  and the manifest `"name"` must match. (Renamed `mozeidon` → `mozeidon_z` in extension **5.0.4**;
  the underscore is mandatory — host names must match `^\w+(\.\w+)*$`, **no hyphens**, so
  `mozeidon-z` is invalid. Changing it again means a coordinated extension + manifest change and an
  AMO re-submit.) The native app itself doesn't reference the host name — it's launched via the
  manifest `"path"` and only uses the socket name.

The **binary filename** (`mozeidon-z-messaging`) and host name are ours; the socket name is frozen.

## Commands

```bash
go build -o mozeidon-z-messaging .   # build
go test ./...                        # test
go vet ./...                         # static check (matches CI)
```

## Release

Bump nothing in code (version comes from the git tag via ldflags). Tag and push:

```bash
git tag -a v1.0.0 -m "mozeidon-z-messaging 1.0.0"
git push origin v1.0.0
```

Needs the `HOMEBREW_TAP_TOKEN` repo secret (PAT with write to `colangelo/homebrew-tap`).

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md). stdout is the native-messaging channel — **never** log to
it; use stderr (`log.Printf`).
