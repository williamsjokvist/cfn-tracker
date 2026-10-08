# Save data and results location settings

Let users pick where the database and the results text files are stored, from the settings page. Defaults stay as they are.

## Today

- `cfn-tracker.db` lives in `os.UserCacheDir()/cfn-tracker`, together with `cfn-tracker.ini` and the Chrome profile.
- `results/*.txt` goes in `results/` relative to the working directory. OBS text sources read from there.

Only the database and results move. The ini has to stay put since it stores the paths, and moving the Chrome profile would log the user out.

## Settings

Add a `[data]` section to `cfn-tracker.ini` with `dir` and `results_dir`. An empty value means the default.

On the settings page, add a "Locations" section with a row for each one. A row shows the current path, plus Change (folder picker), Open folder, and Reset to default.

## Save data

`SaveDataDirectory(dir)`:
- refuses while tracking
- copies the DB with `VACUUM INTO`
- runs migrations on the copy and saves the setting
- keeps the old file as `cfn-tracker.db.bak`

If the target already has a `cfn-tracker.db`, ask whether to use it or replace it.

The new location applies on restart. The SQL storage is shared between handlers, and swapping it live isn't worth it.

If the configured folder is missing at startup (e.g. an unplugged drive), show an error. Don't create an empty DB, since that looks like the data was lost.

## Results

`SaveResultsDirectory(dir)` copies the current `.txt` files over and switches `txt.Storage` to the new folder right away. That needs a mutex, since `SaveMatch` runs on the tracking goroutine. Leave the old files in place, and remind the user to repoint their OBS sources.

If the folder is missing at startup, log a warning and fall back to the default.

## Notes

- New i18n and error keys go in all three locales.
- Regenerate bindings with `mise run bind`. Don't edit them by hand.
- Tests:
  - the path resolvers
  - copy, `.bak` and refusal-while-tracking for the DB move
  - the results switch under `-race`
  - both startup fallbacks
- Open question: should the Chrome profile get its own subfolder?
