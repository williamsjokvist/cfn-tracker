# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Street Fighter 6 tracking.
- Concept of *Sessions* which contain all the matches played during that tracking session.
- Output for OBS Browser Source, with customizeable CSS themes.
- Authentication is now persisted across app restarts.

### Changes
- Tracking errors now render an error notification, others render an error page.
- The config for the GUI is now stored in a *json* file instead of localStorage.
- Sessions, matches and CFN names are now stored in an *SQL* database instead of a big *json* file.
- Match log now distinguishes sessions from matches.

### Fixed
- A Windows XP error message is now displayed when the app crashes on start.

### Removed
- Exporting match log to Excel sheet.
- No longer output tracking state to a *results.json* file.
- Ability to track battle lounge and casual matches **(SFV)**.

## [3.0.0] - 2023-02-22

### Added
- GUI frontend.
- A product page with direct downloads.
- More data associated with the matches are logged.
- The last session can now be resumed.
- A notification is now displayed when there's an update available.
- The match log is now available through the app, with filters for date, opponent, league and character.
- The match log can now be exported to an Excel sheet.
- In addition to text files, the app now also outputs the live data to a JSON file: *results.json*.

### Fixed
- Fixed an issue where the tracking could not be initialized.

Thanks to [Sheldon](https://www.twitch.tv/SheldonTwitching) for testing this release and for the encouragement and support!

Expect the next major release when Street Fighter 6 comes out, until then it's patch work 💯

## [2.0.0] - 2023-01-12

This is the rewrite of the CFN Tracker I wrote almost 4 years ago. It contains many improvements and I hope to add even more features to it in time!

### Added
- Rewrote the app in Go + [Rod](https://github.com/go-rod/rod), since [PhantomJS](https://github.com/ariya/phantomjs) and [CasperJS](https://github.com/casperjs/casperjs) are no longer maintained

### Fixed
- App is now shipped as a single binary instead of having to download PhantomJS + CasperJS.