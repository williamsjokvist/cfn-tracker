# Plan: move the browser source port from `.env` to Settings

## Goal

Remove the `BROWSER_SOURCE_PORT` env config and make the browser source port a
setting in the Settings page, stored in `cfn-tracker.ini` like the other
settings. Changing it takes effect immediately, without restarting the app.

## Current state

- `app/pkg/config/config.go`: `BuildConfig.BrowserSourcePort` is read with
  `envconfig` from `BROWSER_SOURCE_PORT` (default `4242`).
- `app/main.go` `init()`: loads `.env` with `godotenv`. Without a `.env` it
  hardcodes `BrowserSourcePort: 4242`. This is the only env var the app reads,
  so `godotenv` and `envconfig` exist only for it.
- `app/pkg/server/server.go`: `Start` registers handlers on
  `http.DefaultServeMux` and calls `http.ListenAndServe` on the port. There is
  no way to stop or restart the server, and calling `Start` twice panics on
  duplicate handler registration.
- `app/gui/src/pages/output.tsx`: the "Copy Browser Source Link" URL hardcodes
  `http://localhost:4242`, so a custom port already produces a broken link.
- `app/pkg/server/static/index.html`: the setup hint hardcodes the same URL.
- `app/example.env`, `README.md` (dev setup step 2) and
  `app/cmd/init_test.go` reference the env var or the field.

## Changes

### 1. Store the port in runtime config

Follow the `LogFile` setting (commit 7adb01c) as the template.

- `app/pkg/model/config.go`: add
  `BrowserSourcePort uint16 \`ini:"browser_source_port" json:"browserSourcePort" default:"4242"\``
  to `GUIConfig`. It lives in `GUIConfig` (like `log_file`) so the frontend
  gets it from `GetGuiConfig` and the existing `ConfigContext`.
- `app/pkg/storage/config/storage.go`:
  - `GetRuntimeConfig`: default `BrowserSourcePort: 4242` alongside
    `LogFile: true`, so existing ini files without the key get 4242.
  - `SaveRuntimeConfig`: write the `browser_source_port` key.

### 2. Make the server restartable

`app/pkg/server/server.go`:

- Build an `http.ServeMux` once in `NewBrowserSourceServer` and register the
  routes on it instead of on `http.DefaultServeMux`.
- Split the match-forwarding goroutine out of `Start` so it runs once.
- Hold the running `*http.Server` behind the existing mutex and add
  `Listen(port uint16) error`:
  1. `net.Listen("tcp", fmt.Sprintf(":%d", port))` first, so a port that's
     in use or not allowed returns an error to the caller instead of only
     being logged.
  2. Shut down the previous server, if any. SSE connections in
     `handleStream` never finish, so pass a short timeout to `Shutdown` and
     call `Close` afterwards, or have `handleStream` also exit when a
     server-scoped context is cancelled.
  3. `go srv.Serve(listener)` with the new server.
- `Start(ctx, port)` takes the port instead of `*config.BuildConfig`.
- Extend `server_test.go`: listening on a new port stops the old one, and a
  port that's in use returns an error and keeps the old server running.

### 3. Command handler

`app/cmd/cmd.go`:

- Give `CommandHandler` a reference to the browser source server (constructor
  argument, wired in `main.go`).
- Add `SaveBrowserSourcePort(port uint16) error`:
  1. Reject ports below 1024 (`ErrInvalidBrowserSourcePort`).
  2. Call `server.Listen(port)`. On error, return
     `ErrSaveBrowserSourcePort` and don't persist anything.
  3. Save the runtime config.
- Add the error keys in `app/pkg/model/error.go` and the matching
  `Localization` fields in `app/pkg/model/i18n.go`.

### 4. Remove the env config

- `app/pkg/config/config.go`: remove `BrowserSourcePort`. `BuildConfig` keeps
  only `AppVersion`.
- `app/main.go`: remove the `godotenv`/`envconfig` block in `init()` and set
  `cfg.AppVersion` from `wails.json` directly. Read the port from
  `runtimeCfg.GUI.BrowserSourcePort` (already loaded at the top of `main()`)
  and pass it to `browserSrcServer.Start`.
- `app/go.mod`: drop `github.com/joho/godotenv` and
  `github.com/kelseyhightower/envconfig` with `go mod tidy`.
- Delete `app/example.env`.
- `README.md`: remove dev setup step 2 about `example.env`.
- `.gitignore`: leave the `.env` entries; they're harmless and protect
  against committing secrets.
- `app/cmd/init_test.go`: drop `BrowserSourcePort` from the test config.

### 5. Frontend

- `app/gui/src/main/config.tsx`: add `browserSourcePort: 4242` to
  `initialConfig`.
- `app/gui/src/pages/settings.tsx`: add a `BrowserSourcePortInput` row next to
  `LogFileToggle`. A number input that saves on blur or Enter (not on every
  keystroke, since every save rebinds the server), shows errors with
  `useErrorPopup`, and resets to the stored value on failure.
- `app/gui/src/pages/output.tsx`: build the link from
  `cfg.browserSourcePort` (via `ConfigContext`) instead of hardcoded `4242`.
- `app/pkg/server/static/index.html`: build the example URL in the setup hint
  from `location.port` instead of hardcoding 4242.
- Locales (`en-GB`, `fr-FR`, `ja-JP`): add `browserSourcePort` for the label
  plus the two error messages.
- Regenerate `app/gui/wailsjs/**` with `mise run bind` (see `AGENTS.md`). Don't
  edit the bindings by hand.

### 6. Docs

- `web/src/content/changelog.md`: note that the browser source port moved from
  `.env` to Settings.

## Migration

Users who set `BROWSER_SOURCE_PORT` in `.env` fall back to 4242 after
updating and have to set the port again in Settings. Only people who build
from source or put a `.env` next to the binary are affected, so a one-off
migration that reads `.env` isn't worth the dependency. Call it out in the
changelog.

## Verification

- `go test ./...` in `app`, including the new server tests.
- Manual, with `mise app:dev`:
  - Fresh ini: the server listens on 4242 and the copied link works in OBS or
    a browser.
  - Change the port in Settings: the old port stops responding, the new one
    serves, a connected overlay reconnects after its URL is updated, and the
    copied link uses the new port.
  - Pick a port that's in use: an error popup appears, the old port keeps
    serving, and the setting is unchanged after reopening the app.
  - Restart the app: it listens on the saved port.

## Out of scope

- Binding to an interface other than all interfaces (`:port`).
- Picking a free port automatically.
