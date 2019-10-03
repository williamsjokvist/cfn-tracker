@echo off

rename config.ini config.bat
call config.bat
rename config.bat config.ini

title CFN Scraper - %profile%
start "CFN Scraper - %profile%" casperjs/bin/casperjs.exe --web-security=no scraper-min.js --profile="%profile%" --steamid="%steamid%" --steampass="%steampass%  --refreshinterval=%refreshinterval% --writejson=%writejson% --jsonlocation=%jsonlocation%