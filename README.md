```
                   __  _             _  _    ___                 
/\_/\ ___   _   _ / _\| |__    __ _ | || |  / _ \ __ _  ___  ___ 
\_ _// _ \ | | | |\ \ | '_ \  / _` || || | / /_)// _` |/ __|/ __|
 / \| (_) || |_| |_\ \| | | || (_| || || |/ ___/| (_| |\__ \\__ \
 \_/ \___/  \__,_|\__/|_| |_| \__,_||_||_|\/     \__,_||___/|___/
 
 //Curator: Bitly Twiser
                                                                 
```
# YouShallPass
- Password Generator
- Utilized as a web UI, base API to query, or command line application.

# Usage:
- Run install script which will perform the following actions:
1. Compile go binary
2. Place binary within /usr/local/bin to be executed by calling the program name. Default name is "ysp"
- Change name?
- Simply edit the Install File and alter the "binName" variable. Run install script after making modification.

## Command Line:
- After installing binary just call the following:
```ysp --help```
- This will display help menu on program usage.
- Help Menu Output: 
```
  -length int
        Length of the password to generate. (default 8)
  -server
        Starts an Rest API to be queried for passwords. Also will generate a web UI running on port 8080
  -special
        Determines if one desires to DISABLE special characters within the generated password. (default true)
  -upper
        Determines if one desires to have upper case characters within the generated password. (default true)
  -wordlist string
        Location of wordlist to utilize for password generation
 ```

## Web UI:
- To run the web us perform the following:
```ysp -server=true```
- This will start to webserver, browse to port 8080 on localhost to view/use web UI.

## API:
- If the webserver option is utilized on can use the application's built in API to obtain a password.
- The port utilized for the API is __8002__
- The following example uses the HTTPIE project.
- Link: https://httpie.org/
- Usage:
```
http -f POST 127.0.0.1:8002/password length=15 uppercase=false specchar=false
```
- NOTE: __Ensure the webserver is running in order to use API!__
