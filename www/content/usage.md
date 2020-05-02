---
title: "Usage"
menu: true
weight: 4
---

You can pass customize some options.

```sh
> web-watcher --help
Web-watcher discord Bot.

Options:
  -interval int
        The watcher interval in minutes (default 60)
  -prefix string
        The discord commands prefix (default "!")
  -ratio float
        Changes detection ratio (default 1)
  -token string
        Discord token

```

## Arguments

`--interval`

The watcher interval in minutes (default 60).

The watcher will check for website changes at this given interval.

`--prefix`

The discord commands prefix (default !).

`--ratio`

The web page changes ratio. Must be between 0.0 and 1.0.

Every x minutes the watcher will fetch the website page and compares it with the previous version. It will check changes
and convert these changes to a ratio. If page are identical, this ratio is equals to 1.0, and it will decrease for every
detected change.

`--token`

The discord bot token. The token can also be passed with the `DISCORD_TOKEN` environment variable.

If you don't know how to generate one, a quick tutorial describes all the steps in the [requirements page](/web-watcher/requirements/).