# bots-fw-telegram

Telegram module for Strongo bots framework

<!-- dev-approach:v1 -->
## Our approach to development

We build with our own tooling:

- **[SpecScore](https://specscore.md)** — specify requirements as `SpecScore.md` artifacts
- **[SpecStudio](https://specscore.studio)** — author & manage specs across their lifecycle
- **[inGitDB](https://ingitdb.com)** — store structured data in Git where applicable
- **[DALgo](https://dalgo.io)** — data access layer for Go
- **[cover100.dev](https://cover100.dev)** — drive toward 100% test coverage
- **[DataTug](https://datatug.io)** — query & explore data
<!-- /dev-approach -->

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

## Registering webhooks in production

The same handler exposes two routes, mounted by the consuming app under whatever `pathPrefix`
it chooses (e.g. `/bot`):

| Method | Path | Purpose |
|---|---|---|
| `GET` | `{prefix}/tg/set-webhook?code=BOT_CODE` | Calls Telegram's `setWebhook` API for the given bot. Registration is idempotent — safe to call again any time you need to re-point or refresh a webhook (e.g. after adding a secret token). |
| `POST` | `{prefix}/tg/hook?id=BOT_CODE` | The live webhook Telegram POSTs updates to. This is what gets registered by the call above. |

`set-webhook` builds the URL it registers with Telegram from the **incoming request's `Host` header**
(`https://<r.Host>{prefix}/tg/hook?id=BOT_CODE`), not from a fixed config value. This matters in production:

- **Call it on whatever host you want Telegram to POST updates to.** If your app sits behind a reverse
  proxy / CDN worker that rewrites the `Host` header before forwarding to the origin (e.g. a Cloudflare
  Worker fronting a Cloud Run service without a native domain mapping), the registered webhook URL will
  reflect the **origin's** hostname, not the public-facing one you called. Check what `Host` your origin
  actually receives before registering through a proxied domain, or call `set-webhook` directly on the
  origin's own hostname if you want predictability.
- Each bot's `BotSettings.WebhookSecretToken` (if configured) is sent as Telegram's `secret_token`, and
  incoming webhook requests are checked against it (`X-Telegram-Bot-Api-Secret-Token` header). Registering
  without one leaves that bot's webhook unauthenticated — anyone who learns/guesses the webhook URL can
  POST forged updates. A warning is logged (not fatal) when this happens.
- Bot tokens are resolved from `<PLATFORM>_BOT_TOKEN_<CODE>` env vars by default (e.g.
  `TELEGRAM_BOT_TOKEN_MYBOT` for a bot registered with code `MyBot`), unless the consuming app
  wires a token explicitly.

## Used by

- sneat-go (private)