# myShrugBot - A Telegram bot to express your indifference
[![build](https://github.com/d-Rickyy-b/myShrugBot/actions/workflows/release_build.yml/badge.svg)](https://github.com/d-Rickyy-b/myShrugBot/actions/workflows/release_build.yml)
[![test](https://github.com/d-Rickyy-b/myShrugBot/actions/workflows/test_push_pr.yml/badge.svg)](https://github.com/d-Rickyy-b/myShrugBot/actions/workflows/test_push_pr.yml)

Quickly send a `¯\_(ツ)_/¯` within a Telegram chat with this simple bot.
You can find the hosted version at [@myShrugBot](https://t.me/myShrugBot) on Telegram.

## Configuration
This bot uses a configuration file. You can create your own using the sample file `config.sample.yml`.

```yaml
bot_token: "<your_bot_token>"

webhook:
  enabled: false
  listen_ip: "127.0.0.1"
  listen_port: 9001
  listen_path: "/"
  url: "example.com/test"
  cert_path: ""
  cert_key_path: ""
```
Just fill in your bot token. If you want to use a webhook, use the webhook options.
The current config only allows for receiving webhooks from a reverse proxy such as nginx.

## Usage
Create your `config.yml` file and put it into the same directory as the binary. 
You can then simply call the binary.

```
Usage of myShrugBot:
  -config string
        Path to config file (default "config.yml")
```

You can also use the parameter `-config` to speficy a path to another config file.
