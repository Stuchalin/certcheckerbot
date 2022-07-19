# certcheckerbot

Simple bot to check certificate expiry date

## Environment configuration
```
BOT_KEY=Secret telegram bot key
DB_PATH=path to sqlite database
DEBUG=true/false (enable or disable debug. default - false)
```

## Available commands
**/help** - bot commands help

**/check [www.checkURL1.com www.checkURL2.com ...]** - check cert info for transferred domain names