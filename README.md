# certcheckerbot

Simple bot to check certificate expiry date

## Environment configuration
```
BOT_KEY=Secret telegram bot key
DB_PATH=path to sqlite database
DEBUG=true/false (enable or disable debug. default - false)
EXPIRY_DAYS=[1,2,3,4,5,6,7,14,30,60,90]
```

## Available commands
**/help** - bot commands help

**/check [www.checkURL1.com www.checkURL2.com ...]** - check certificate on URL. Use spaces to check few domains

**/set_hour [hour in 24 format 0..23]** - set a notification hour for messages about expired domains. For example: "/set_hour 9". Notification hour for default - 0.

**/set_tz [-11..14]** - set a timezone for messages about expired domains. For example: "/set_tz 3". Timezone for default - 0.

**/domains** - get added domains

**/add_domain [domain_name]** - add domain for schedule checks. For example: "/add_domain google.com"

**/remove_domain [domain_name]** - removes domain for schedule checks. For example: "/remove_domain google.com"

## v0.3
* Work all base commands