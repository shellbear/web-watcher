---
title: "Requirements"
menu: true
weight: 2
---

First of all, you have to create a new discord bot, generate a token and add the bot to the server of your choice.

If you don't know how to so, here are some few steps:

Go to the [discord developer portal](https://discordapp.com/developers/applications) and create a new application with the name of your choice:

![App creation](/web-watcher/app-creation.png)

Click on the `Bot` tab on the left and create a new bot:

![Bot creation](/web-watcher/bot-creation.png)

Then go to the `OAuth2` tab, check the `Bot` scope and the `Send messages` permission:

![Bot permissions](/web-watcher/bot-permissions.png)

**⚠️ Recent Update**

Due to Discord API changes, you must enable `Message Content Intent` privilege for the bot to work correctly.

![Bot privilege](/web-watcher/bot-privilege.png)

Then, copy the generated URL in the middle of the screen, open it and add the bot to the server of your choice.

Congrats, you added the bot to your discord server!

Now all you have to do is to obtain your secret token. This is the token that will be used by web-watcher to connect to
discord with the `--token` option or `DISCORD_TOKEN` environment variable. For more details check the [usage page](/web-watcher/usage).

Return to the `Bot` tab and click on `Click to Reveal Token`:

![Bot token](/web-watcher/bot-token.png)

Copy it for the next steps and make sure to keep this token secret!
