# Botio

[![Go Report Card](https://goreportcard.com/badge/github.com/danielkvist/botio)](https://goreportcard.com/report/github.com/danielkvist/botio)
[![Docker Pulls](https://img.shields.io/docker/pulls/danielkvist/botio.svg?maxAge=604800)](https://hub.docker.com/r/danielkvist/botio/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

Botio is a simple and opinionated CLI to create and manage easily bots for differents platforms."

## Platforms supported

- Telegram.
- Discord.

## Example

First, we create a server from which we're going to manage the commands to our bot:

```bash
botio server --db ./data/botio.db --col commands --http :9090 --key mysupersecretkey
```

We add a command:

```bash
botio add --command start --response Hello --url :9090 --key mysupersecretkey
```

We can check that the command has been created successfully with ```print```:

```bash
botio print --command start --url :9090 --key mysupersecretkey
```

Or see a list of all our commands with ```list```:

```bash
botio list --url :9090 --key mysupersecretkey
```

Now, we can start our Telegram's bot:

```bash
botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey
```

> If you have doubts about how to get a Telegram token for your bot you can check this [link](https://core.telegram.org/bots#botfather).

To check all the available commands use the ```help``` flag:

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
botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey
botio print --command start --url :9090 --key mysupersecretkey

Available Commands:
  add         Adds a new command
  delete      Deletes the specified command
  discord     Initializes a Discord bot
  help        Help about any command
  list        Prints a list with all the commands
  print       Prints the specified command and his response
  server      Starts a server to manage the commands with simple HTTP methods.
  telegram    Initializes a Telegram bot
  update      Updates an existing command (or adds it if not exists)

Flags:
  -h, --help   help for botio

Use "botio [command] --help" for more information about a command.
```

### server

```text
$ botio server --help
Starts a server to manage the commands with simple HTTP methods.

Usage:
  botio server [flags]

Examples:
botio server --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey

Flags:
  -c, --collection string   collection used to store commands (default "commands")
  -d, --database string     database path (default "./commands.db")
  -h, --help                help for server
      --http string         port for HTTP connections (default ":80")
      --https string        port for HTTPS connections (default ":443")
  -k, --key string          authentication key
      --sslcert string      ssl certification file
      --sslkey string       ssl key file
```

Example:

```bash
botio server --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey
```

> The database used is based on BoltDB. You can read more about it [here](https://github.com/etcd-io/bbolt).

### telegram

```text
$ botio telegram --help
Initializes a Telegram bot

Usage:
  botio telegram [flags]

Examples:
botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey

Flags:
  -h, --help           help for telegram
  -k, --key string     authentication key
  -t, --token string   telegram's token
  -u, --url string     botio's server URL
```

Example:

```bash
botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey
```

### discord

```text
$ botio discord --help
Initializes a Discord bot

Usage:
  botio discord [flags]

Examples:
botio discord --token <discord-token> --url :9090 --key mysupersecretkey

Flags:
  -h, --help           help for discord
  -k, --key string     authentication key
  -t, --token string   discord's token
  -u, --url string     botio's server URL
```

Example:

```bash
botio discord --token <discord-token> --url :9090 --key mysupersecretkey
```

### add

```text
$ botio add --help
Adds a new command

Usage:
  botio add [flags]

Examples:
botio add --command start --response Hello --url :9090 --key mysupersecretkey

Flags:
  -c, --command string    command to add
  -h, --help              help for add
  -k, --key string        authentication key
  -r, --response string   command's response
  -u, --url string        botio's server url
```

Example:

```bash
bodio add --command start --response Hello --url :9090 --key mysupersecretpassword
```

### print

```text
$ botio print --help
Prints the specified command and his response

Usage:
  botio print [flags]

Examples:
botio print --command start --url :9090 --key mysupersecretkey

Flags:
  -c, --command string   command to print
  -h, --help             help for print
  -k, --key string       authentication key
  -u, --url string       botio's server URL
```

Example:

```bash
botio print --command start --url :9090 --key mysupersecretkey
```

### list

```text
$ botio list --help
Prints a list with all the commands

Usage:
  botio list [flags]

Examples:
botio list --url :9090 --key mysupersecretkey

Flags:
  -h, --help         help for list
  -k, --key string   authentication key
  -u, --url string   botio's server URL
```

Example:

```bash
botio list --url :9090 --key mysupersecretkey
```

### update

```text
$ botio update --help
Updates an existing command (or adds it if not exists)

Usage:
  botio update [flags]

Examples:
botio update --command start --response Hi --url :9090 --key mysupersecretkey

Flags:
  -c, --command string    command to update
  -h, --help              help for update
  -k, --key string        authentication key
  -r, --response string   command's new response
  -u, --url string        botio's server url
```

Example:

```text
botio update --command start --response Hi --url :9090 --key mysupersecretkey
```

### delete

```text
$ botio delete --help
Deletes the specified command

Usage:
  botio delete [flags]

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
botio delete --command start --url :9090 --key mysupersecretkey
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

### Backup

> Backup is an special endpoint that will send to the client a backup of the database.

```text
http://<url>:<port>/api/backup
```

```text
http://localhost:9090/api/backup
```

## ToDo

- [ ] Better logging
- [ ] Docker Compose
- [ ] Web Interface
- [ ] Alternative databases like PostgreSQL
- [ ] Support for Facebook Messenger bots
- [ ] Support for Slack bots
- [ ] Support for Skype bots
