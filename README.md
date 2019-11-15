# Botio

[![Go Report Card](https://goreportcard.com/badge/github.com/danielkvist/botio)](https://goreportcard.com/report/github.com/danielkvist/botio)
[![CircleCI](https://circleci.com/gh/danielkvist/botio.svg?style=svg)](https://circleci.com/gh/danielkvist/botio)
[![Docker Pulls](https://img.shields.io/docker/pulls/danielkvist/botio.svg?maxAge=604800)](https://hub.docker.com/r/danielkvist/botio/)
[![LICENSE](https://img.shields.io/github/license/danielkvist/botio)](https://github.com/danielkvist/botio/blob/master/LICENSE)
[![Issues](https://img.shields.io/github/issues/danielkvist/botio)](https://github.com/danielkvist/botio/issues)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)
[![GoDoc](https://godoc.org/github.com/danielkvist/botio?status.svg)](https://godoc.org/github.com/danielkvist/botio)

Botio is a CLI to create and manage easily chatbots for different platforms with the possibility of using differents databases.

> NOTE: Due to how filepaths work if you're a Windows user is recommended to use a shell like Git Bash.

## Support

### Platforms

- Telegram.
- Discord.

#### Work is in progress to add support for:

- Slack.
- Facebook Messenger.
- Alexa.
- Google Assistant.

### Databases

- BoltDB
- PostgreSQL

#### Work is in progress to add support for:

- MongoDB.
- Redis.
- SQLite.

## Install

### Go

```bash
go get -u github.com/danielkvist/botio
```

### Docker

```bash
docker image pull danielkvist/botio
```

## Development

### Clone

```bash
git clone https://github.com/danielkvist/botio.git botio
```

### Setup

```bash
cd ./botio
```

And

```bash
make setup
```

## CLI

Botio provides a simple CLI to manage your server, your chatbots and their commands.

```text
$ botio --help
Botio is a CLI to create and manage easily chatbots for different platforms such as Telegram or Discord.  
It also let's you use different databases to manage their available commands wuch as BoltDB or PostgreSQL.

Botio is a project in development so use it with caution!

Usage:
  botio [command]

Examples:
botio server bolt --database ./data/commands.db --collection commands
botio bot --platform telegram --token <telegram-token>
botio client add --command start --response Hi

Available Commands:
  bot         Starts a chatbot for the specified platform.
  client      Client provides subcommands to manage your commands.
  help        Help about any command
  server      Server provides subcommands to initialize a server with differents databases.

Flags:
  -h, --help   help for botio

Use "botio [command] --help" for more information about a command.
```

### Server

The `server` subcommand handles the initialization of a Botio server with the specified database.

```text
$ botio server --help
Server provides subcommands to initialize a server with differents databases.

Usage:
  botio server [command]

Available Commands:
  bolt        Starts a Botio server with BoltDB.
  postgres    Starts a Botio server with PostgreSQL.

Flags:
  -h, --help   help for server

Use "botio server [command] --help" for more information about a command.
```

If you want to initialize a Botio's server with BoltDB:

```bash
botio server bolt
```

Or with PostgreSQL:

```bash
botio server postgres
```

> IMPORTANT: Due to how PostgreSQL works you will need to have created the database before trying to connect Botio to it.

### Client

In the future you will have the option to manage your chabots commands with an HTTP client. For the moment you can use the `client` subcommand.

```text
$ botio client --help
Client provides subcommands to manage your commands.

Usage:
  botio client [command]

Available Commands:
  add         Adds a new command.
  delete      Deletes the requested command
  list        List all the commands.
  print       Prints the requested command.
  update      Updates the requested command or adds it if don't exists.  

Flags:
  -h, --help   help for client

Use "botio client [command] --help" for more information about a command.
```

Each subcommand provides and example so feel free to check each one by one.

### Bot


The `bot` subcommand handles the initialization of a chatbot for a specified platform.

```text
$ botio bot --help
Starts a chatbot for the specified platform.

Usage:
  botio bot [flags]

Examples:
botio bot --platform telegram --token <telegram-token>

Flags:
      --addr string       botio's gRPC server address (default ":9091")
      --goroutines int    number of goroutines (default 10)
  -h, --help              help for bot
      --platform string   platform (discord or telegram)
      --sslca string      ssl client certification file
      --sslcrt string     ssl certification file
      --sslkey string     ssl certification key file
      --token string      bot's token
```

If for example you want to initialize a chatbot for Telegram:

```bash
botio bot --platform telegram --token <telegram-token>
```

> Please, check the documentation provided by the differents plaforms about how to get a token for a chatbot.

## gRPC HTTP endpoint

Botio uses the Google's gRPC gateway to provide an HTTP client. For the moment is work in pogress and it doesn't work.
Once it is finished all the corresponding documentation will be here. If you want to help PRs they are welcome.

## Other things that need to improve

You can secure with TLS your server or not. To do this you simply have to leave the flags `--sslca`, `--sslcrt` and `--sslkey` empty. The same goes for the client and the chabot's client.

I know that right now managing the configuration and generation of certificates to get a secure server or client is a bit clumsy and I'm actively working to solve it.
