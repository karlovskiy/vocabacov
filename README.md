# Vocabacov

Telegram bot to enlarge your vocabulary.

## Environment variables
| Name               | Mandatory | Default value | Description                                   |
|--------------------|-----------|---------------|-----------------------------------------------|
| VOCABACOV_TOKEN    | yes       |               | Telegram bot token.                           |
| VOCABACOV_CHANNELS | yes       |               | Comma delimited ids of telegram bot channels. |
| VOCABACOV_DEBUG    | no        | false         | Debug logging enabled/disabled.               |
| VOCABACOV_TIMEOUT  | no        | 30            | Request timeout to wait for an update.        |

## Run

To run bot with telegram token `YOUR_TELEGRAM_BOT_TOKEN`, channel `36484` and mapping volume database file
`/db/vocabacov.db` to you host's current directory:
```shell
touch $(pwd)/vocabacov.db

docker run -d \
--name vocabacov \
-e VOCABACOV_TOKEN=YOUR_TELEGRAM_BOT_TOKEN \
-e VOCABACOV_CHANNELS=36484 \
-v $(pwd)/vocabacov.db:/db/vocabacov.db \
ghcr.io/karlovskiy/vocabacov:latest
```

## Telegram bot

To create telegram bot follow this [guide](https://core.telegram.org/bots).

## Commands

`Vocabacov` bot accepts commands in the format `/lang phrase`, 
where `lang` is ISO-3166 Alpha-2 code and phrase is one or more words.

Examples: `/en hello world`, `/es ola chica`

This bot saves submitted words in the `sqlite3` database.
```shell
$ sqlite3
SQLite version 3.38.2 2022-03-26 13:51:10
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
sqlite> .open vocabacov.db
sqlite> select * from phrases;
1|en|hello world
2|en|all your base are belong to us
```

## TODO
Automatic creation of [anki](https://apps.ankiweb.net/) collections from `sqlite` database 
with help of awesome [genanki](https://github.com/kerrickstaley/genanki) library.