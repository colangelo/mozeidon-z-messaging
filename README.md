# mozeidon-z-messaging

The native-messaging host for the **Mozeidon-Z** stack — a tiny Go proxy that lets the
[Mozeidon-Z browser extension](https://addons.mozilla.org/firefox/addon/mozeidon-z/) exchange
commands and responses with the [Mozeidon-Z CLI](https://github.com/colangelo/mozeidon-z) over a
local IPC socket.

> Hard fork of [`egovelox/mozeidon-native-app`](https://github.com/egovelox/mozeidon-native-app).
> The **binary** is renamed `mozeidon-z-messaging`; the native-messaging **host name** (`mozeidon`)
> and the **IPC socket** (`mozeidon_native_app`) are unchanged, so it is a drop-in for the
> Mozeidon-Z extension with no AMO change. See [ARCHITECTURE.md](ARCHITECTURE.md) for how it works.

## Install (macOS / Linux)

```bash
brew install colangelo/tap/mozeidon-z-messaging
```

This is also pulled automatically as a dependency of `brew install colangelo/tap/mozeidon-z`.

## Configure native messaging (Firefox, macOS)

Create `~/Library/Application Support/Mozilla/NativeMessagingHosts/mozeidon.json`:

```json
{
  "name": "mozeidon",
  "description": "Mozeidon native messaging host",
  "path": "/opt/homebrew/bin/mozeidon-z-messaging",
  "type": "stdio",
  "allowed_extensions": ["mozeidon-z@a-layer.io"]
}
```

(`just setup-native-messaging` in the `mozeidon-z` repo writes this for you.) Restart Firefox.

## Usage

It is launched by the browser, not run directly. For diagnostics:

```bash
mozeidon-z-messaging --version
mozeidon-z-messaging --help
```

## Build from source

```bash
go build -o mozeidon-z-messaging .   # needs Go 1.26+
```

## Releases

A `v*` git tag triggers GitHub Actions (`.github/workflows/release.yml`) → a matrix build of
`mozeidon-z-messaging-<os>-<arch>` for darwin + linux (amd64/arm64), cosign keyless signing, a public
GitHub Release, and an `update-homebrew` job that bumps `Formula/mozeidon-z-messaging.rb` in
`colangelo/homebrew-tap`. (Same pattern as the `mozeidon-z` CLI.)

## License

MIT. Originally based on `egovelox/mozeidon-native-app`.

