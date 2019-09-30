# Botio

[![Go Report Card](https://goreportcard.com/badge/github.com/danielkvist/botio)](https://goreportcard.com/report/github.com/danielkvist/botio)
[![CircleCI](https://circleci.com/gh/danielkvist/botio.svg?style=svg)](https://circleci.com/gh/danielkvist/botio)
[![Docker Pulls](https://img.shields.io/docker/pulls/danielkvist/botio.svg?maxAge=604800)](https://hub.docker.com/r/danielkvist/botio/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

Botio is a simple and opinionated CLI to create and manage easily bots for differents platforms.

## Supported platforms

- Telegram.
- Discord.

## Example

First, we create a server from which we're going to manage the commands to our bot:

```bash
botio server --db ./data/botio.db --col commands --http :9090 --key mysupersecretkey
```

We add a command:

```bash
botio client add --command start --response Hello --url :9090 --key mysupersecretkey
```

We can check that the command has been created successfully with `print`:

```bash
botio client print --command start --url :9090 --key mysupersecretkey
```

Or see a list of all our commands with `list`:

```bash
botio client list --url :9090 --key mysupersecretkey
```

Now, we can start a Telegram's bot:

```bash
botio bot --platform telegram --token <telegram-token> --url :9090 --key mysupersecretkey
```

To check all the available commands use the `help` flag:

```bash
botio help
```

And that's it, now all that's left is to add or edit commands according to our needs.

## Install

### Go

```bash
go install github.com/danielkvist/botio
```

### Docker

```bash
docker image pull danielkvist/botio
```

## CLI

```text
$ botio help
Botio is a simple and opinionated CLI to create and manage easily bots for differents platforms.

Usage:
  botio [command]

Examples:
botio server --database ./data/commands.db --collection commands --http :9090 --key mysupersecretkey
botio bot --platform telegram --token <telegram-token> --url :9090 --key mysupersecretkey
botio client print --command start --url :9090 --key mysupersecretkey

Available Commands:
  bot         Initializes a bot for a supported platform (telegram and discord for the moment)
  client      Client contains some subcommands to manage your bot's commands
  help        Help about any command
  server      Starts a server to manage the commands with simple HTTP methods

Flags:
  -h, --help   help for botio

Use "botio [command] --help" for more information about a command.
```

### server

```text
$ botio server --help
server contains some subcommands to initialize a server with different databases

Usage:
  botio server [flags]
  botio server [command]

Available Commands:
  bolt        Starts a server with a BoltDB database to manage your commands with HTTP methods
  postgres    Starts a server with that connects to a PostgreSQL database to manage your commands with HTTP methods

Flags:
  -h, --help   help for server

Use "botio server [command] --help" for more information about a command.
```

#### bolt

```text
$ botio server bolt --help
Starts a server with a BoltDB database to manage your commands with HTTP methods

Usage:
  botio server bolt [flags]

Examples:
botio server bolt --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey

Flags:
      --collection string   collection used to store commands (default "commands")
      --database string     database path (default "./commands.db")
  -h, --help                help for bolt
      --http string         port for HTTP connections (default ":80")
      --https string        port for HTTPS connections (default ":443")
      --key string          authentication key
      --sslcert string      ssl certification file
      --sslkey string       ssl certification key file
```

Example:

```bash
botio server bolt --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey
```

#### postgres

```text
$ botio server postgres --help
Starts a server with that connects to a PostgreSQL database to manage your commands with HTTP methods

Usage:
  botio server postgres [flags]

Examples:
botio server postgres --user postgres --password toor --database botio --table commands --key mysupersecretkey

Flags:
      --database string   PostgreSQL database name (default "botio")
  -h, --help              help for postgres
      --host string       host of the PostgreSQL database (default "postgres")
      --http string       port for HTTP connections (default ":80")
      --https string      port for HTTPS connections (default ":443")
      --key string        authentication key
      --password string   password for the user of the PostgreSQL database
      --port string       port of the PostgreSQL database host (default "5432")
      --sslcert string    ssl certification file
      --sslkey string     ssl certification key file
      --table string      table of the PostgreSQL database
      --user string       user of the PostgreSQL database
```

Example:

```bash
botio server postgres --user postgres --password toor --database botio --table commands --key mysupersecretkey
```

### bot

```text
$ botio bot --help
Initializes a bot for a supported platform (telegram and discord for the moment)

Usage:
  botio bot [flags]

Examples:
botio bot --platform telegram --token <telegram-token> --url :9090 --key mysupersecretkey

Flags:
  -h, --help              help for bot
  -k, --key string        authentication key
  -p, --platform string   platform (discord or telegram)
  -t, --token string      bot's token
  -u, --url string        botio's server URL
```

Example:

```bash
botio bot --platform telegram --token <telegram-token> --url :9090 --key mysupersecretkey
```

### client

```text
$ botio client --help
Client contains some subcommands to manage your bot's commands

Usage:
  botio client
  botio client [command]

Available Commands:
  add         Adds a new command
  delete      Deletes the specified command
  list        Prints a list with all the commands
  print       Prints the specified command and his response
  update      Updates an existing command (or adds it if not exists)

Flags:
  -h, --help   help for client

Use "botio client [command] --help" for more information about a command.
```

#### add

```text
$ botio client add --help
Adds a new command

Usage:
  botio client add [flags]

Examples:
botio client add --command start --response Hello --url :9090 --key mysupersecretkey

Flags:
  -c, --command string    command to add
  -h, --help              help for add
  -k, --key string        authentication key
  -r, --response string   command's response
  -u, --url string        botio's server url
```

Example:

```bash
bodio client add --command start --response Hello --url :9090 --key mysupersecretpassword
```

#### print

```text
$ botio client print --help
Prints the specified command and his response

Usage:
  botio client print [flags]

Examples:
botio client print --command start --url :9090 --key mysupersecretkey

Flags:
  -c, --command string   command to print
  -h, --help             help for print
  -k, --key string       authentication key
  -u, --url string       botio's server URL
```

Example:

```bash
botio client print --command start --url :9090 --key mysupersecretkey
```

### list

```text
$ botio client list --help
Prints a list with all the commands

Usage:
  botio client list [flags]

Examples:
botio client list --url :9090 --key mysupersecretkey

Flags:
  -h, --help         help for list
  -k, --key string   authentication key
  -u, --url string   botio's server URL
```

Example:

```bash
botio client list --url :9090 --key mysupersecretkey
```

### update

```text
$ botio client update --help
Updates an existing command (or adds it if not exists)

Usage:
  botio client update [flags]

Examples:
botio client update --command start --response Hi --url :9090 --key mysupersecretkey

Flags:
  -c, --command string    command to update
  -h, --help              help for update
  -k, --key string        authentication key
  -r, --response string   command's new response
  -u, --url string        botio's server url
```

Example:

```text
botio client update --command start --response Hi --url :9090 --key mysupersecretkey
```

### delete

```text
$ botio client delete --help
Deletes the specified command

Usage:
  botio client delete [flags]

Examples:
botio delete --command start --url :9090 --key mysupersecretkey

Flags:
  -c, --command string   command to delete
  -h, --help             help for delete
  -k, --key string       authentication key
  -u, --url string       botio's server url
```

Example:

```bash
botio client delete --command start --url :9090 --key mysupersecretkey
```

## API Endpoints

### GET

#### GET Commands

```text
http://<url>:<port>/api/commands
```

```text
http://localhost:9090/api/commands
```

#### GET Command

```text
http://<url>:<port>/api/commands/<command>
```

```text
http://localhost:9090/api/commands/start
```

### POST

```text
http://<url>:<port>/api/commands
```

```text
http://localhost:9090/api/commands
```

### UPDATE

```text
http://<url>:<port>/api/commands/<command>
```

```text
http://localhost:9090/api/commands/start
```

### DELETE

```text
http://<url>:<port>/api/commands/<command>
```

```text
http://localhost:9090/api/commands/start
```

## ToDo

- [ ] Docker Compose
- [ ] Web Interface
- [ ] Support for Facebook Messenger bots
- [ ] Support for Slack bots
- [ ] Support for Skype bots
