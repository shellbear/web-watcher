# Web-watcher

A small discord bot which aims to monitor and send notifications on website changes.

## Installation

If you want to self host this bot, you have to first create a new Discord application and bot from the [developer portal](https://discordapp.com/developers/applications/).
You can follow this [tutorial](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token) to achieve this step.

With go CLI:
```bash
go get github.com/ShellBear/web-watcher
```

## Usage

```bash
export DISCORD_TOKEN=YOUR_DISCORD_TOKEN
web-watcher

# or

web-watcher --token YOUR_DISCORD_TOKEN
```

```bash
web-watcher --help

Web-watcher discord Bot.

Options:
  -delay int
        Watch delay in minutes (default 60)
  -token string
        Discord token
```

By default the watch interval for every website is 1 Hour but you can easily change this with the `delay` parameter followed by the time interval in minutes.

```bash
# Set watch interval to 10 minutes
web-watcher --token YOUR_DISCORD_TOKEN --delay 10
```

## Commands

#### !watch [URL]

Add an URL to the watchlist.

#### !unwatch [URL]

Remove an URL from the watchlist.

#### !watchlist

Get the complete watchlist.

## Built With

- [Gorm](https://github.com/jinzhu/gorm) - The fantastic ORM library for Golang
- [DiscordGo](https://github.com/bwmarrin/discordgo) - Go bindings for Discord

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details