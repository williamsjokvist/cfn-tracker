/*
	cfn profile scraper
	author: greensoap
*/

var casper = require("casper").create({
	pageSettings: {loadImages:  false, loadPlugins: false},
	viewportSize: {width: 1, height: 1}
}),

fs = require('fs'),
    
/*Configuration*/
cfn = casper.cli.raw.get("profile"),
refreshInterval = casper.cli.raw.get("refreshinterval"),
writeAJSON = casper.cli.raw.get("writejson"), 
jsonLoc = casper.cli.raw.get("jsonlocation"),

modes = [
    {
        num: 0,
        type: "rank",
        wins: 0,
        losses: 0,
        ratio: "0%"
    },
    {
        num: 1,
        type: "casual",
        wins: 0,
        losses: 0,
        ratio: "0%"
    },
    {
        num: 2,
        type: "lounge",
        wins: 0,
        losses: 0,
        ratio: "0%"
    }
];

var updated = 0, /* = 1 if match stats have been updated*/
    updateCount = [0,0,0], /*each modes update counter*/
    modeArr = [[], [], []]; /*array which contains the three modes stat arrays*/

/**
 * Select W/L nodelist and calls updateStats with the new array. Reloads after it's done.
 */
function fetchStats() {   
    var data = [0], arr = [0];
    casper.waitForSelector('.' + modes[2].type + ' .win .on', function () {
        for (var i = 0; i < modes.length; i++){           
            data[i] = this.evaluate(function (modeType){
                return document.querySelector('.'+ modeType +" .win ul").innerHTML;
            }, modes[i].type);

            arr[i] = data[i].split('</'); /*Converts the nodelist into an array*/

            updateStats(arr[i], modes[i]);

            reload();
        }    
    });
}

var myRank = 0,
    myLP = 0, /*parseInt()*/
    netLP = 0,
    LPArray = [];

/**
 * Fetches Rank and League Points
 */
function fetchRank() {
    var rankSelector = ".playerInfo>dl:last-child>dd";
    
    casper.waitForSelector(rankSelector, function (){
        myRank = this.evaluate(function (selector){
            return document.querySelector(selector).innerText;
        }, rankSelector);

        var lp = this.evaluate(function (){
            return document.querySelector(".leagueInfo>dl:last-child>dd").innerText;
        }); 
        
        myLP = parseInt(lp, 10);/*remove the "LP"*/
        
        /*Net LP*/
        LPArray.push(myLP);
        if (LPArray.length > 1 && (LPArray[LPArray.length-1] !== LPArray[LPArray.length-2])){
            netLP = LPArray[LPArray.length-1] - LPArray[0];
        
            if (netLP > 0)/*positive*/
                netLP = "+" + netLP;
        }
                
        this.echo("\n\trank: " + myRank + "\t\tLP: " + myLP + " (" + getLeague(myLP) + ")\tNet LP: " + netLP);
    });
}
  


/**
 * Checks if W/L counters needs updating, and does it.
 * @param {matchArr} the mode's match array
 * @param {mode} mode which stats are to be checked for updates
 */
function updateStats(matchArr, mode){
    for (var i = 0; i < matchArr.length; i++){
        if (matchArr[i] !== modeArr[mode.num][i]){/*if new win or loss*/    
            var today = new Date(),
            time = today.getHours() + ":" + today.getMinutes();
            
            /*Ádd to W/L counters*/
            if (matchArr[19].indexOf('on') !== -1)
                mode.wins++;
             else 
                mode.losses++;
            
            /*Ignore first fetch*/
            if (updateCount[mode.num] === 0 && ((mode.wins === 0 && mode.losses === 1) || (mode.wins === 1 && mode.losses === 0))){
                mode.wins = 0;
                mode.losses = 0;                  
            } else if ((mode.losses > 0) || (mode.wins > 0)){
                mode.ratio = Math.floor((mode.wins / (mode.wins+mode.losses)) * 100) + "%";
                var matchMeta = "\r\n("+ time + ") " + mode.type + " match #" + updateCount[mode.num],
                    matchStat = "\r\n\tWins: " + mode.wins + "\t\t\tLosses: " + mode.losses + "\t\t\tWin Ratio: " + mode.ratio,
                    cMatchStat = "Wins: " + mode.wins + ", Losses: " + mode.losses + ", Win Ratio: " + mode.ratio;
                
                casper.echo(matchMeta);
                writeLog("\t" + matchMeta + "\t" + cMatchStat + "\t");
                if (mode.type === "rank"){
                    fetchVersus();
                    fetchRank();
                    casper.wait(500, function(){this.echo(matchStat);});
                } else
                    casper.echo(matchStat);                    
            }
            
            modeArr[mode.num] = matchArr.slice(); /*Updates the array*/
            writeTxt(mode.type, mode.wins, mode.losses, mode.ratio);
            
            if (writeAJSON === "true")
                writeJSON(0);
            
            updated = 1;     
            updateCount[mode.num]++;
        }
        else {
            updated = 0;
        }
    }
}

/*
* Fetch the opponents name and LP as well as the characters used
*/
function fetchVersus(){
    casper.waitForSelector(".recentbox .fighter", function (){
        var opName = this.evaluate(function (){
            return document.querySelectorAll(".recentbox>.fId>dd")[0].innerText;
        }),
        opLp = this.evaluate(function (){
            return document.querySelectorAll(".recentbox>.league>dd")[0].innerText;
        }),
        opChar = this.evaluate(function (){
            return (document.querySelectorAll(".recentbox>.fav>dd")[0].innerText);
        }),
        myChar = this.evaluate(function (){
            return (document.querySelectorAll(".leagueWrap dl>dd>a")[0].innerText);
        }),
        string = cfn + " (" + globalizeChar(myChar) + ") " + "(" + getLeague(myLP) + ") vs. " + opName + " (" + globalizeChar(opChar) + ") " + "(" + getLeague(opLp) + ")";
        
        this.echo("\t" + string);
        writeLog("\n" + string);
    });
}

function globalizeChar(string){
    if (string === "VEGA") return "Boxer";
    else if (string === "M. BISON") return "Claw";/*they never fixed this*/
    else if (string === "BALROG") return "Dictator";    
    else return string.charAt(0) + string.slice(1).toLowerCase();/*return capitalized string*/
}

function getLeague (lp){
    var league = "";
    if (lp >= 300000) league = "Warlord";
    else if (lp >= 100000) league = "Ultimate Grand Master";
    else if (lp >= 35000) league = "Grand Master";
    else if (lp >= 30000) league = "Master";
    else if (lp >= 25000) league = "Ultra Diamond";
    else if (lp >= 20000) league = "Super Diamond";
    else if (lp >= 14000) league = "Diamond";
    else if (lp >= 12000) league = "Ultra Platinum";
    else if (lp >= 10000) league = "Super Platinum";
    else if (lp >= 7500) league = "Platinum";
    else if (lp >= 6500) league = "Ultra Gold";
    else if (lp >= 4500) league = "Super Gold";
    else if (lp >= 4000) league = "Gold";
    else if (lp >= 3500) league = "Ultra Silver";
    else if (lp >= 3000) league = "Super Silver";
    else if (lp >= 2000) league = "Silver";
    else if (lp >= 1500) league = "Ultra Bronze";
    else if (lp >= 1000) league = "Super Bronze";
    else if (lp >= 500) league = "Bronze";
    else if (lp < 500) league = "Rookie";
    
    return league;
}

/**
 * Refreshes the page
 */
function reload(){
    casper.wait(refreshInterval, function(){
        this.reload(function(){
            fetchStats();            
        });
    });
}

/**
 * Skip these unnecessary requests
 */
casper.options.onResourceRequested = function(casper, requestData, request) {
  var skip = [
	'https://www.gstatic.com/charts/loader.js',
	'https://www.gstatic.com/charts/46.1/loader.js',
	'https://www.gstatic.com/charts/46.1/js/jsapi_compiled_corechart_module.js',
	'https://www.gstatic.com/charts/46.1/js/jsapi_compiled_ui_module.js',
	'https://www.gstatic.com/charts/46.1/js/jsapi_compiled_default_module.js',
	'https://www.gstatic.com/charts/46.1/js/jsapi_compiled_format_module.js'
  ];

  skip.forEach(function(needle) {
    if (requestData.url.indexOf(needle) > 0) 
        request.abort();    
  })
};

/**
 * Writes .txt files for each mode as well as LP and rank
 * @param {mode} the mode which wins/losses/ratio is written to
 * @param {wins} the amount of wins
 * @param {losses} the amount of losses
 * @param {ratio} the win ratio
 */
function writeTxt (mode, wins, losses, ratio) {
    fs.write(fs.pathJoin('result/' + mode, 'wins.txt'), wins);
    fs.write(fs.pathJoin('result/' + mode, 'losses.txt'), losses);
    fs.write(fs.pathJoin('result/' + mode, 'ratio.txt'), ratio);
    fs.write(fs.pathJoin('result/rank.txt'), myRank);
    fs.write(fs.pathJoin('result/LP.txt'), myLP);
    fs.write(fs.pathJoin('result/netLP.txt'), netLP);
}
/**
 * Generate log file / append to log file
 * @param {data} data to append
 */
function writeLog(data){
    fs.write("LOGFiLE.log", data, 'a');
}
	
/**
 * Writes a JSON and sends it to the jsonLocation
 * @param {clear} clear wins, losses and ratio, only at launch
 */
function writeJSON (clear){
    var arr, subArr = [];
    
    if (clear === 0){
        for (var i = 0; i < modes.length;i++){
            subArr[i] = '"' + modes[i].type + '":{'
            + '"wins":' + modes[i].wins + ','
            + '"losses":' + modes[i].losses +  ',' 
            + '"ratio":"' + modes[i].ratio + '"}';
        }   
    } else {
        for (var i = 0; i < modes.length;i++){
            subArr[i] = '"' + modes[i].type + '":{"wins":' + 0 + ',losses":' + 0 +  ',ratio":"' + 0 + '"}';
        }
    }
    
    arr = '{"account":{"rank":' + myRank + ',"lp":' + myLP + ',"net":"' + netLP + '"},' + subArr[0] + ", " + subArr[1] + ", " + subArr[2] + "}";

    var xhr = new XMLHttpRequest();

    xhr.open('PUT', jsonLoc);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        if (xhr.status === 200 && xhr.responseText !== arr)
            casper.echo('Something went wrong.  It is currently' + xhr.responseText);          
        else if (xhr.status !== 200)
            casper.echo('Request failed.  Returned status of ' + xhr.status);  
        
    };
    xhr.send(arr);
}

/**
 * Run for your life!
 */
function run (){
    var today = new Date(),
    date = today.getFullYear()+'-'+(today.getMonth()+1)+'-'+today.getDate();
    casper.echo(date);
    writeLog(date);
    casper.start("https://game.capcom.com/cfn/sfv/gate/steam?rpnt=https://game.capcom.com/cfn/sfv/profile/" + cfn, function () {
        /*Start and agree to terms*/
        this.echo("Action: Accepting CFN Terms");
        this.click('input[value="Agree"]');
    });

    casper.waitForSelector("form input[name='username']", function () {
        /*Fill form and log in*/
        this.echo("Action: Logging in to " + this.getTitle());
        this.fillSelectors('form#loginForm', {
            'input[name = username ]': casper.cli.raw.get("steamid"),
            'input[name = password ]': casper.cli.raw.get("steampass")
        }, true);

        this.click('input#imageLogin');
        this.echo("Success");
        this.wait(3000);
    });

    fetchRank();
    fetchStats();
    
    if (writeAJSON === "true")
        writeJSON(1);
    
    casper.run();
}

if ((refreshInterval < 5000 || refreshInterval === null) || cfn === null){
	casper.echo('Your settings are misconfigured');
	casper.exit();
} else 
    run();