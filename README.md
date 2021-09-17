# rr-parser
raid.report telegram bot (player stats in chat)

To use that, you need to do some simple steps:
1. Create your own telegram bot https://core.telegram.org/bots#creating-a-new-bot
2. Create your OpenSSL certificates https://www.openssl.org/ (command in certs/certs.go)
3. Create PostgreSQL database and execute db.sql file
4. Get your Bungie API key https://www.bungie.net/developer
5. Set up config.json
6. Compile and run :)

Bot commands:
1. **/rr** - check player stats (format: **/rr raid nickname**)
2. **/reg** - register your Destiny 2 nickname to your tg profile (format: **/reg destiny2_nickname**)
3. **/upd** - update your registered Destiny 2 nickname (format: **/upd destiny2_nickname**)
4. **/my** - check your stats (format: **/my raid**). Only for registered users.
5. **/lists** - raids abbreviation
6. **/help** - commands and raids abbreviation
