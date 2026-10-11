# Save data location and backups

Let users pick where the database is stored, and back it up and restore it, from the settings page. The default location stays as it is.

## Today

- `cfn-tracker.db` lives in `os.UserCacheDir()/cfn-tracker`, together with `cfn-tracker.ini` and the Chrome profile.

Only the database moves. The ini has to stay put since it stores the path, and moving the Chrome profile would log the user out.

The results text files are out of scope. They keep their own location and are handled separately.

## Settings

Add a `[data]` section to `cfn-tracker.ini` with `dir`. An empty value means the default.

On the settings page, add a "Save data" section. It shows the current path, plus Change (folder picker), Open folder and Reset to default, and Back up and Import backup.

## Location

`SaveDataDirectory(dir)`:
- refuses while tracking
- copies the DB with `VACUUM INTO` and saves the setting
- leaves the old DB where it is, since it's still open until the restart (and Windows can't rename an open file)

If the target already has a `cfn-tracker.db`, ask whether to use it or replace it. Using it checks it the same way as an import. Migrations run on the next start, like for any DB.

A move copies the DB the app would open next, so a pending move or import isn't lost by moving again.

The new location applies on restart. The SQL storage is shared between handlers, and swapping it live isn't worth it.

If the configured folder is missing at startup (e.g. an unplugged drive), show an error. Don't create an empty DB, since that looks like the data was lost.

## Dialogs

The CSV export already gets a save dialog: `Exporter` takes a `SaveFileDialogFn`, set in `OnStartup` with `runtime.SaveFileDialog`, and cancelling isn't an error. Backup, import and Change all follow that pattern. Add an open-file and a directory dialog func next to it, so the handlers stay testable without Wails.

## Backup

`BackupData()` opens a save dialog with `cfn-tracker.db.bak` as the default filename, and writes the backup wherever the user picks with `VACUUM INTO`. `VACUUM INTO` gives a consistent copy even while the DB is open, so this works while tracking. `VACUUM INTO` won't write over an existing file, so write to a temp file in the target folder and rename it over the target.

## Import

`ImportBackup()` opens a file dialog (`.bak` and `.db`) so the user can pick a backup from anywhere, and restores it over `cfn-tracker.db`.

- refuses while tracking
- asks for confirmation, since it overwrites the current data
- checks that the file is a cfn-tracker DB (opens it and finds the `migrations` table) and that its schema isn't newer than the app's
- copies it to `cfn-tracker.db.import` next to the DB

At startup, before opening the DB, a `cfn-tracker.db.import` replaces `cfn-tracker.db` and migrations run as usual. Like the location change, it applies on restart.

## Notes

- New i18n and error keys go in all locales. `TestLocalesMatchLocalizationStruct` fails otherwise.
- Regenerate bindings with `mise run bind`. Don't edit them by hand.
- Tests:
  - the path resolver
  - copy, `.bak` and refusal-while-tracking for the move
  - backup while the DB is open, and over an existing file
  - cancelling a dialog is a no-op
  - import: validation, refusal-while-tracking, and the swap at startup
  - the missing-folder startup error
