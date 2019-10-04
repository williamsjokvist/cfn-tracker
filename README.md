# CFN Scraper
This tool scrapes any CFN profile for its wins, losses, rank, LP as well as your Net LP.

### Who is this for? 
This is primarily targeted toward streamers who want automatically updated win/loss counters. It can also be useful if you wanna keep track of your LP gain during a session.

### Configuration and setting it up

* Get [PhantomJS](https://phantomjs.org/download.html) and put the .exe in the master directory. 
* Get [CasperJS](http://casperjs.org/), change the name of the folder inside the .zip file to "casperjs" and put it in the directory as well. It should now look like [this](showcase/folder-structure.jpg? "folder-structure").
* Set up a Steam account with Steam Guard **disabled** – this is because logging in is required to view CFN profiles. :unamused: 
  The password cannot contain a '%'.
* In the _config.ini_ you will need to specify the Steam ID and Password for the account, as well as the target CFN account.
* You have the options to output the result to individual .txt files as well as a JSON in case you wanna do some fancy stuff ^^
* After you've configured it to your liking, launch the .bat file. .txt files will be generated in a folder called "result". Use them as sources for your labels and they'll update after every match.

### Example

![stream](showcase/screenie.gif?raw=true "streamshowcase")

### Result

This is how the batch file will display the data:
![screenshot](showcase/batchscreenshot.GIF?raw=true "screenshot")


Here's what the JSON looks like:
```
{
   "account":{
      "rank":0,
      "lp":0,
      "net":"0"
   },
   "rank":{
      "wins":0,
      "losses":0,
      "ratio":"0%"
   },
   "casual":{
      "wins":0,
      "losses":0,
      "ratio":"0%"
   },
   "lounge":{
      "wins":0,
      "losses":0,
      "ratio":"0%"
   }
}
```
