<p align="center">
  <img alt="Gopher" src=".github/images/gopher.png" height="140" />
  <h3 align="center">web-watcher</h3>
  <p align="center">A small Discord bot which aims to alert you on website changes.</p>
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/shellbear/web-watcher" alt="Go Report Card">
    <img src="https://goreportcard.com/badge/github.com/shellbear/web-watcher" />
  </a>
  <a href="https://github.com/shellbear/web-watcher/actions?query=workflow%3Alint" alt="Pipeline status">
    <img src="https://github.com/shellbear/web-watcher/workflows/lint/badge.svg" />
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/shellbear/web-watcher" alt="Go version" />
  <a href="https://opensource.org/licenses/MIT" alt="Go version">
    <img src="https://img.shields.io/badge/license-MIT-brightgreen.svg" />
  </a>
</p>

---

## Installation

If you want to self host this bot, you have to first create a new Discord application and bot from the [developer portal](https://discordapp.com/developers/applications/).
You can follow this [tutorial](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token) to achieve this step.

Grab you discord token and then run the bot.

## Usage

```bash
> web-watcher --help
Web-watcher discord Bot.

Options:
  -delay int
        Watch delay in minutes (default 60)
  -prefix string
        The discord commands prefix (default "!")
  -ratio float
        Changes detection ratio (default 1)
  -token string
        Discord token
```

By default, the watch interval for every website is 1 hour, but you can easily change this with the `delay` parameter followed by the time interval in minutes.

```bash
# Set watch interval to 10 minutes
web-watcher --token YOUR_DISCORD_TOKEN --delay 10
```

## Discord commands

#### !watch [URL]

Add a URL to the watchlist.

#### !unwatch [URL]

Remove a URL from the watchlist.

#### !watchlist

Get the complete watchlist.

## Deploy

With Docker:
```bash
docker build -t web-watcher .
```

Run:
```bash
docker run -e DISCORD_TOKEN=YOUR_DISCORD_TOKEN web-watcher
```

## Limitations

Web-watcher analyzes the static HTML page structure. Changes detection, is based on page structure and tags, text modifications are not yet supported.
It's only supporting static HTML pages at the moment.

## Built With

- [go-difflib](https://github.com/pmezard/go-difflib) - Partial port of Python difflib package to Go 
- [xxhash](https://github.com/cespare/xxhash) - Go implementation of the 64-bit xxHash algorithm (XXH64)
- [Gorm](https://github.com/jinzhu/gorm) - The fantastic ORM library for Golang
- [DiscordGo](https://github.com/bwmarrin/discordgo) - Go bindings for Discord

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details