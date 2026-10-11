# FAQ

## The app shows up as a virus in Windows Defender.
This is a false positive, caused by the app not being signed by a certificate authority. App certificates are expensive - [around $500/yr](https://order.digicert.com/step1/code_signing) - and no way am I paying that. The [source code](https://github.com/williamsjokvist/cfn-tracker) is public, so you or a programming friend can check it out and make sure no tomfoolery is going on 😌

## Can I use this for my own without streaming?
Yes! Quite a few people have also told me that peeking at their stats, during or after they play has helped them improve.

## The Mac version says it's damaged or can't be opened.
The Mac app isn't signed with an Apple Developer ID, so macOS blocks it the first time. Download [the Mac version](https://github.com/williamsjokvist/cfn-tracker/releases/latest/download/cfn-tracker-macos-universal.zip), unzip it and move it to Applications. Then either try to open it and click "Open Anyway" under System Settings → Privacy & Security, or run this in Terminal:

```sh
xattr -dr com.apple.quarantine "/Applications/CFN Tracker.app"
```

## What is the future of CFN Tracker?
Adding more statistics overall, especially related to sessions, character and player matchups over time

## How can I report bugs?
Open an [issue](https://github.com/williamsjokvist/cfn-tracker/issues) on GitHub, or DM on Twitter
