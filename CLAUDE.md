# CLAUDE.md

Guidance for Claude Code working in this repo.

## What this is

`mozeidon-z-messaging` — the native-messaging host (browser ⇄ CLI IPC bridge) for the Mozeidon-Z
stack. ~230 lines of Go. Hard fork of `egovelox/mozeidon-native-app` (remote `upstream`).

## Frozen identifiers — DO NOT CHANGE

- Native-messaging **host name** `"mozeidon"` — the shipped AMO extension calls
  `connectNative("mozeidon")`. Changing it forces an AMO re-submit.
- IPC **socket** base name `mozeidon_native_app` (generated form
  `mozeidon_native_app_<pid>_<profileId8>`) — contract with the `mozeidon-z` CLI.

Only the **binary filename** (`mozeidon-z-messaging`) is ours to rename.

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
