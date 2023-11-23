<img src="build/appicon.png" height="150px" align="right"/>

# CFN Tracker
This tool tracks any Street Fighter 6 or V CFN account's live match statistics. 

## Who is this for? 
This is primarily targeted toward streamers who want automatically updated win/loss counters. It can also be useful if you wanna keep track of your LP gain and other statistics during a session.

## Support
If you find the tool useful for your streams or otherwise, consider [buying me a coffee ☕](https://ko-fi.com/greensoap)

If you encounter any bugs, @ me on Twitter [@greensoap_](https://twitter.com/GreenSoap_)

## Getting Started
* Download the app for either Mac or Windows from [the download page](https://williamsjokvist.github.io/cfn-tracker/).
* Extract the app wherever you want and pick either SFV or SF6
* Input your CFN when prompted and start
* If all is well, tracking has now started. Go to the `Output` screen, pick a theme and mark the stats you want to display and click `Copy Browser Source Link`
* Open OBS and add a Browser Source, paste the link as the source.

Text files will also be generated in a `results` folder, which are updated after every match. For memory efficiency you can use these instead of the `Browser Source` in OBS. 

## Features

### Real-Time Tracking
![tracking](https://raw.githubusercontent.com/williamsjokvist/cfn-tracker/github-pages/public/tracking.png "tracking-example")

### View the Match Log
![log](https://raw.githubusercontent.com/williamsjokvist/cfn-tracker/github-pages/public/log.png "log-example")

### Display Stats on Stream
Display the stats either via an OBS Browser Source or by importing simple text files. 
![browser-src](https://raw.githubusercontent.com/williamsjokvist/cfn-tracker/github-pages/public/output.png "game-pick")

### Track Multiple Games
![game](https://raw.githubusercontent.com/williamsjokvist/cfn-tracker/github-pages/public/game.png "game-pick")

## Creating your own Browser Source theme with CSS
Create a CSS file in the `themes` directory and it will be available as an option in the theme selector. In the CSS there are 4 CSS selectors you can use: 
- `stat-list`: the container for all of the stats
- `stat-list::part(list-item)`: the container for one stat
- `stat-list::part(list-title)`: the stat title e.g. "LP Gain", "Win Rate" ...
- `stat-list::part(list-value)`: the value of the stat

Easiest is to base a new theme on `nord.css`

## Streamer showcase
### [OneStepLayered](https://twitch.tv/OneStepLayered)
![stream](https://raw.githubusercontent.com/williamsjokvist/cfn-tracker/github-pages/public/osl.gif "OneStepLayered")
### [SheldonTwitching](https://twitch.tv/SheldonTwitching)
![stream](https://raw.githubusercontent.com/williamsjokvist/cfn-tracker/github-pages/public/sheldon.jpg "SheldonTwitching")