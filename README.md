# bots-fw-telegram

Telegram module for Strongo bots framework

## Structure & key concepts

The [`tgWebhookHandler`](tg_webhook_handler.go) struct is implementing `botsfw.WebhookHandler` interface
and is an entry point for all incoming requests from Telegram. To create it you need to call
[`NewTgWebhookHandler()`](tg_webhook_handler.go) function.

### Registering records maker

```go
```

## Setting up dev environment for Telegram bots development

## Tunneling to local development environment

To expose local server to the Internet we use [ngrok](https://ngrok.com/).

```shell
ngrok http 4300
```

Make sure that you have started local GAE server & Firestore emulators - follow instructions
from [README.md](README.md).

After `ngrok` started you will see something like:

```shell
Forwarding    https://****-***-**.ngrok-free.app -> http://localhost:4300
```

You would need to register the forwarding URL for bot you are testing with Telegram by calling this url:

`https://****-***-**.ngrok-free.app/bot/tg/set-webhook?code=BOT_CODE`

where `****-***-**` is the forwarding URL from `ngrok` output and `BOT_CODE` is the code of the bot you are testing.

The bot will be registered using secret tokens that you should set using environment variables:

```shell
TELEGRAM_BOT_TOKEN_<BOT_CODE>=<TELEGRAM_BOT_TOKEN>
```

You can create a personal bot for testing purposes using [BotFather](https://t.me/botfather).

The bot with the given code should be registered in your app and the value is CASE SENSITIVE.