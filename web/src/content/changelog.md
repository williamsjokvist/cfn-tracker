# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [5.0.2](https://github.com/williamsjokvist/cfn-tracker/tree/v5.0.2)
## Fixed
- App crashing on inputting a Tekken ID without a recent match
- Allow inputting dashes in Tekken ID, like shown in-game
- Character name mapping for Lidia, Clive, Anna and Fahkumram

## [5.0.1](https://github.com/williamsjokvist/cfn-tracker/tree/v5.0.1)
## Fixed
- App crashing on bad Tekken ID input
- Adjusted text message on Tekken ID input to no longer say "CFN Code"

## [5.0.0](https://github.com/williamsjokvist/cfn-tracker/tree/v5.0.0)

## Added
- Tekken 8 support
- Themes for the app UI
- Sessions calendar page
- Force update tracking button

## Changed
- Settings have been moved from the side bar to a separate page
- Updates are now automatic instead of prompted in the app

## Fixed
- All stats are now logged in text files instead of just a couple

## Removed
- Support for SFV

## [4.1.0](https://github.com/williamsjokvist/cfn-tracker/tree/v4.1.0)

### Fixed
- Browser source now instantly renders stats instead of waiting until a match has been played
- Fixed issue where custom themes weren't being applied
- Improved theme selection in output settings
- Adjusted theme CSS structure

## [4.0.1](https://github.com/williamsjokvist/cfn-tracker/tree/v4.0.1)

### Fixed
- Fixed issue where wins and losses were not being counted

## [4.0.0](https://github.com/williamsjokvist/cfn-tracker/tree/v4.0.0)

### Added
- Street Fighter 6 tracking.
- Concept of *Sessions* which contain all the matches played during that tracking session.
- Output for OBS Browser Source, with ability to make custom themes in CSS.
- Authentication is now persisted across app starts.

### Changed
- Tracking errors now render an error notification, others render an error page.
- The config for the GUI is now stored in a *json* file instead of localStorage.
- Sessions, matches and CFN names are now stored in an *SQL* database instead of a big *json* file.
- *Match Log* renamed to *Sessions*, and lists sessions instead of all matches.
- GUI State is controlled via a state machine.

### Fixed
- LP and MR Gains are now counted by character **(SF6)**
- In cases of app crashing on start, a Windows XP error message is displayed.

### Removed
- Exporting match log to Excel sheet.
- No longer output tracking state to a *results.json* file.
- Ability to track battle lounge and casual matches **(SFV)**.

## [3.0.0](https://github.com/williamsjokvist/cfn-tracker/tree/v3.0.0) - 2023-02-22

Thanks to [Sheldon](https://www.twitch.tv/SheldonTwitching) for testing this release and for the encouragement and support!

Expect the next major release when Street Fighter 6 comes out, until then it's patch work 💯

### Added
- GUI frontend.
- Translations for *French* and *Japanese*.
- A product page with direct downloads.
- More data associated with the matches are logged.
- The last session can now be resumed.
- A notification is displayed when there's an update available.
- Match log, with filters for date, opponent, league and character, and Excel export.
- In addition to text files, the tracking data is also outputted to a JSON file: *results.json*

### Fixed
- Fixed an issue where the tracking could not be initialized.

## [2.0.0](https://github.com/williamsjokvist/cfn-tracker/tree/v2.0.0) - 2023-01-12

This is the rewrite of the CFN Tracker I wrote almost 4 years ago. It contains many improvements and I hope to add even more features to it in time!

### Added
- Rewrote the app in Go + [Rod](https://github.com/go-rod/rod), since [PhantomJS](https://github.com/ariya/phantomjs) and [CasperJS](https://github.com/casperjs/casperjs) are no longer maintained

### Changed
- App is now shipped as a single binary instead of having to download PhantomJS + CasperJS.

## [1.0.0](https://github.com/williamsjokvist/cfn-tracker/tree/458774bf59df5854b7ba6365a0f0b3cfc74bc52f) - 2019-09-07
