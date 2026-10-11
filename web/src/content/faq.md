# FAQ

## The app shows up as a virus in Windows Defender.
This is a false positive, caused by the app not being signed by a certificate authority. App certificates are expensive - [around $500/yr](https://order.digicert.com/step1/code_signing) - and no way am I paying that. The [source code](https://github.com/williamsjokvist/cfn-tracker) is public, so you or a programming friend can check it out and make sure no tomfoolery is going on 😌

## The Mac version says it's damaged or can't be opened.
The Mac app isn't signed with an Apple Developer ID, so macOS blocks it the first time. Download [the Mac version](https://github.com/williamsjokvist/cfn-tracker/releases/latest/download/cfn-tracker-macos.zip), unzip it and move it to Applications. Then either try to open it and click "Open Anyway" under System Settings → Privacy & Security, or run this in Terminal:

```sh
xattr -dr com.apple.quarantine "/Applications/CFN Tracker.app"
```

## How can I report bugs?
Use *Report issue* in the app's sidebar. It posts the report publicly on GitHub together with today's app log, which makes the bug much easier to track down. You can also open an [issue](https://github.com/williamsjokvist/cfn-tracker/issues) on GitHub yourself, or DM on Twitter.
