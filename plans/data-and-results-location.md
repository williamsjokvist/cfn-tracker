# Save data and results location settings

Let the user choose where CFN Tracker keeps its save data and its results text files, from a new section on the settings page. The defaults stay where they are today, so nothing changes for users who never touch the settings.

## Current state

| What | Location | Where it's decided |
| --- | --- | --- |
| `cfn-tracker.db`: matches, sessions, users | `os.UserCacheDir()/cfn-tracker` | `getDataSource()` in `pkg/storage/sql/storage.go` |
| `cfn-tracker.ini`: settings | `os.UserCacheDir()/cfn-tracker` | `getAppDataDir()` in `pkg/storage/config/storage.go` |
| Chrome profile with the buckler session | `os.UserCacheDir()/cfn-tracker`, the same folder rather than a subfolder | `UserDataDir()` in `pkg/browser/browser.go` |
| `results/*.txt`: one file per match field, read by OBS text sources | `results/`, relative to the working directory | `NewStorage()` in `pkg/storage/txt` |

`os.UserCacheDir()` is `~/Library/Caches` on macOS and `%LocalAppData%` on Windows.

## Scope

Two new settings:

- **Save data location**, for `cfn-tracker.db`. The default stays `os.UserCacheDir()/cfn-tracker`.
- **Results location**, for the `results/*.txt` files. The default stays `results/` in the working directory. Users' OBS text sources already point there, so changing the default would break them.

These don't move:

- `cfn-tracker.ini`, because it stores the chosen locations and has to be found first.
- The Chrome profile. It's a cache, and moving it would log the user out.
- `cfn-tracker.log`.

## Design

### Storage

- Add a `[data]` section to `cfn-tracker.ini` with `dir` and `results_dir` keys, read through `model.RuntimeConfig`. An empty key means the default.
- Replace the hardcoded path in `getDataSource()` with a resolver that takes the configured directory and falls back to the default.
- Give `txt.Storage` a directory that is set from config at startup instead of the hardcoded `results`.

### Changing the save data location

Add `CommandHandler.SaveDataDirectory(dir string) error`. It refuses while tracking is running, then:

1. Checks that `dir` exists and is writable.
2. If `dir` already has a `cfn-tracker.db`, returns a specific error so the UI can ask whether to use that file or replace it.
3. Copies the current database with `VACUUM INTO '<dir>/cfn-tracker.db'`. That gives a consistent copy while the connection is open, so it's safer than copying the file.
4. Opens the copy and runs the migrations as a check, then saves the new `dir` to the ini.
5. Keeps the old file and renames it to `cfn-tracker.db.bak`, so nothing is lost if the new location turns out to be wrong.

The change takes effect after a restart, and the UI says so. The same `*sql.Storage` is shared by `CommandHandler` and `TrackingHandler`, and swapping it live isn't worth the risk.

`CommandHandler.ResetDataDirectory()` does the same steps toward the default location.

### Changing the results location

Add `CommandHandler.SaveResultsDirectory(dir string) error`:

1. Checks that `dir` exists and is writable.
2. Copies the current `*.txt` files over, so OBS shows the latest values right away once it's repointed.
3. Saves `results_dir` to the ini and calls a new `txt.Storage.SetDirectory(dir)`. The setter is mutex-guarded, since `SaveMatch` runs from the tracking goroutine.

This takes effect immediately, with no restart. The old files are left where they were, in case OBS still points at them.

The UI reminds the user to update their OBS text sources to the new folder.

`CommandHandler.ResetResultsDirectory()` goes back to the default. `OpenResultsDirectory` opens whichever folder is configured.

### Startup

If the configured save data directory is missing, for example an unplugged drive, don't create a new empty database there or at the default. An empty database would look like the user's data was lost. Instead, show the error window with the path, telling the user to reconnect the drive or reset the location.

If the configured results directory is missing, log a warning and fall back to the default. Those files are only a live output, so nothing is lost.

### Settings page

Add a "Locations" section under the existing toggles, with one row each for "Save data" and "Results". Each row has:

- the current path, in a monospace font that can be truncated
- **Change…**, which opens Wails' `runtime.OpenDirectoryDialog` and calls the matching save method
- **Open folder**
- **Reset to default**, shown only when a custom location is set

Changing the save data location shows a notice that the app needs to be restarted. Changing the results location shows the reminder about OBS sources.

### Strings and bindings

- New translation keys in all three locales and in `model.Localization`. They cover the section and row titles, the buttons, the restart notice, the OBS reminder and the "a save file already exists here" prompt.
- New error keys: `errSaveDataDirectory`, `errDataDirectoryMissing` and `errSaveResultsDirectory`.
- Regenerate `app/gui/wailsjs` with `mise run bind` after adding the handler methods. Don't edit it by hand (see `AGENTS.md`).

## Tests

- The path resolvers: an empty config gives the defaults, and configured directories are used as-is.
- `SaveDataDirectory` against temp directories:
  - the copy contains the same rows
  - the ini is updated
  - the old file is kept as `.bak`
  - it's refused while tracking
  - a target with an existing database returns the specific error
- `SaveResultsDirectory`:
  - the current files are copied
  - matches saved afterwards are written to the new folder
  - `SetDirectory` is safe under `-race` while `SaveMatch` runs concurrently
- Startup:
  - a missing save data directory returns an error instead of creating a database
  - a missing results directory falls back to the default

## Open questions

- Should the Chrome profile move into its own subfolder, so it doesn't share a folder with the save data?
