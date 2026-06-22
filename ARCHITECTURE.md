# Architecture — mozeidon-z-messaging

```
Mozeidon-Z browser extension
        ▲  │   native messaging: 4-byte LE length-prefixed JSON over the host's stdin/stdout
        │  ▼
  mozeidon-z-messaging   (this binary — launched by the browser)
        ▲  │   IPC: Unix socket  mozeidon_native_app_<pid>_<profileId8>
        │  ▼
  Mozeidon-Z CLI / Raycast
```

## Two protocols, one proxy

The host bridges two channels:

1. **Browser ↔ host — native messaging.** The browser starts this binary and speaks the
   [native-messaging protocol](https://developer.chrome.com/docs/extensions/develop/concepts/native-messaging):
   each message is a 4-byte little-endian length prefix followed by that many bytes of JSON, over
   the host's **stdin/stdout**. Because stdout *is* the protocol, the host must never print anything
   else there — all logging goes to **stderr** (`log.Printf`).
2. **Host ↔ CLI — IPC.** The host runs a `james-barrow/golang-ipc` server on a Unix socket. The CLI
   connects, sends an `{command, args}` message, and reads streamed responses.

`webBrowserProxy()` is the loop: read an IPC message → forward to the browser
(`PostMessage(os.Stdout, …)`) → read browser responses (`OnMessage(os.Stdin, …)`) and relay each
back over IPC until the `{"data":"end"}` terminator (`isEndOfStream`).

## Registration & multi-profile

On startup the browser sends a first **registration** message (`models.RegistrationInfo`:
browser name/engine/version, `profileId`, rank, name, aliases, user agent, timestamp). The host:

1. Builds a `NativeAppProfile` (`models.GetNativeAppProfile`) with a **per-instance** socket name
   `mozeidon_native_app_<pid>_<profileId8>` and filename `<pid>_<profileId8>.json`.
   (Guarded: a `profileId` shorter than 8 chars is rejected, not panicked on.)
2. Writes that profile JSON into `$UserConfigDir/mozeidon_profiles/`.
3. Starts the IPC server on the per-instance socket.

This lets several browsers/profiles run concurrently, each with its own host instance + socket.

## The 3-way contract

The registration/profile schema is shared across three components — change one, change all three:

| Leg | Component | File |
|---|---|---|
| sends registration | extension | `firefox-addon/src/services/registration.ts` |
| writes profile + socket name | **this host** | `models/*.go` |
| reads profiles, dials the socket | CLI | `mozeidon-z` `cli/profiles/profiles.go`, `cli/core/app.go` |

The CLI also keeps a legacy fallback to the fixed socket `mozeidon_native_app`.

## Lifecycle

On `SIGTERM`/`SIGINT` (the browser closing the host) or any error exit, the host removes its profile
file (`signal.Notify` + `defer os.Remove`), so stale profiles don't accumulate. (Signal-based
unregister does not fire on Windows — which is why we don't ship Windows builds.)

## Security notes

- The IPC socket is created with default (owner-only) permissions (`UnmaskPermissions: false`);
  the host and CLI run as the same user.
- `golang-ipc`'s "encryption" is a homegrown handshake, not audited crypto. The trust model is
  localhost / single-user.
