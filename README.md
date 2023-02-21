<img src="build/appicon.png" height="150px" align="right"/>

# CFN Tracker
This tool tracks any Street Fighter V CFN account's live match statistics. 

## Who is this for? 
This is primarily targeted toward streamers who want automatically updated win/loss counters. It can also be useful if you wanna keep track of your LP gain and other statistics during a session.

## Support

If you find the tool useful for your streams or otherwise, consider [buying me a coffee ☕](https://ko-fi.com/greensoap)

## Usage

* Download the app for either Mac or Windows from [the download page](https://greensoap.github.io/cfn-tracker/).
* Launch it and wait for it to finish initializing.
* Input your CFN when prompted and start.

Text files will be generated in a `results` folder, which are updated after every match. Use them as sources for your labels in e.g. OBS.

All of the logged matches are listed in the `Match Log` tab, and can be filtered by character, date, league and opponent.

The app has translations for `English`, `French` and `Japanese`.

The app will start with a notification if an update has released. You can also follow me on [Twitter](https://twitter.com/GreenSoap_) for update tweets.

## Screenshots

![tracking](showcase/tracking.jpg?raw=true "tracking-example")
![log](showcase/log.jpg?raw=true "log-example")

## Coming updates
* Optional web output, so you can add it as a browser source in OBS
* Themes for said web output, with easy customization of layout, font, font size and colors
* Street Fighter 6 support

Any ideas for features and improvements are welcome, feel free to open an [issue](https://github.com/GreenSoap/cfn-tracker/issues)

If you wish to contribute, you can open a [pull request](https://github.com/GreenSoap/cfn-tracker/pulls) 

## Tools

The app is created with React for the [frontend](https://github.com/GreenSoap/cfn-tracker/tree/master/frontend) and Go for the [backend](https://github.com/GreenSoap/cfn-tracker/tree/master/frontend), using the [Wails](https://github.com/wailsapp/wails) framework.

## Building from source

To build from source, run `make gui` in the root directory. 

Wails is required to build from source: [installation instructions](https://wails.io/docs/gettingstarted/installation).

The only required dependencies for Wails are:
* Go 1.18+
* Node 15+

An .env file following [example.env](https://github.com/GreenSoap/cfn-tracker/blob/master/example.env) is also required.

## Streamers using the app

If you use this app for your stream, DM me a screenshot on [Twitter](https://twitter.com/GreenSoap_) and I will add it here!
### [OneStepLayered](https://twitch.tv/OneStepLayered)
![stream](showcase/stream-example-2.gif?raw=true "stream-example")