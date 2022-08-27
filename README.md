# srakabot2

## Configuring

Bot congiration goes in yaml file of the next format

    telegram:
        api_token: YOUR:TELEGRAM_BOT_APIKEY_HERE
    couchdb:
        host: http://your.couchdb.host:5984/
        user: username
        password: password
        database: database_name

put this as .srakabot.yml to the home directory or use

``` -config /path/to/yaml ``` command line argument to specify location
