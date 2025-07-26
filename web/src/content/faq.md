# FAQ

## The app shows up as a virus in Windows Defender.
This is a false positive, caused by the app not being signed by a certificate authority. App certificates are expensive - [around $500/yr](https://order.digicert.com/step1/code_signing) - and no way am I paying that. The [source code](https://github.com/williamsjokvist/cfn-tracker) is public, so you or a programming friend can check it out and make sure no tomfoolery is going on 😌

## Can I use this for my own without streaming?
Yes! Quite a few people have also told me that peeking at their stats, during or after they play has helped them improve.

## Why was the Mac version removed?
There are many issues releasing apps on Mac, so I spared myself the headache of supporting it. You can still download it [here](https://github.com/williamsjokvist/cfn-tracker/releases/download/v5.0.0/cfn-tracker-darwin-arm64.zip) and try to get it to work, or build from source which works fine.

## What is the future of CFN Tracker?
Adding more statistics overall, especially related to sessions, character and player matchups over time

## How can I report bugs?
Open an [issue](https://github.com/williamsjokvist/cfn-tracker/issues) on GitHub, or DM on Twitter
