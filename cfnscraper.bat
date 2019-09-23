@echo off

rename config.ini config.bat
call config.bat
rename config.bat config.ini

color 0C
title CFN Scraper - %profile%
casperjs --web-security=no scraper-min.js --mode="%mode%" --profile="%profile%" --steamid="%steamid%" --steampass="%steampass%  --refreshinterval=%refreshinterval% --writejson=%writejson% --jsonlocation=%jsonlocation%
pause