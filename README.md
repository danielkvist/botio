# Botio

[![Go Report Card](https://goreportcard.com/badge/github.com/danielkvist/botio)](https://goreportcard.com/report/github.com/danielkvist/botio)
[![Docker Pulls](https://img.shields.io/docker/pulls/danielkvist/botio.svg?maxAge=604800)](https://hub.docker.com/r/danielkvist/botio/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

Botio is an opinionated CLI to create and manage easily simple bots for differents platforms.

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
Simple CLI tool to create and manage easily bots for different platforms.

Usage:
  botio [command]

Examples:
botio server --db ./data/commands.db --col commands --http :9090 --key mysupersecretkey
botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey
botio print --command start --url :9090 --key mysupersecretkey

Available Commands:
  add         Adds a new command with a response to the botio's server
  delete      Deletes the specified botio's command from the botio's server
  discord     Initializes a Discord bot that extracts the commands from the botio's server
  help        Help about any command
  list        Prints a list with all the botio's commands
  print       Prints the specified botio's command with his response
  server      Starts a botio's server to manage the botio's commands with simple HTTP methods.
  telegram    Initializes a Telegram bot that extracts the commands from the botio's server.
  update      Updates an existing command (or adds it if not exists) with a response on the botio's server

Flags:
  -h, --help   help for botio

Use "botio [command] --help" for more information about a command.
```

### server

```text
$ botio server --help
Starts a botio's server to manage the botio's commands with simple HTTP methods.

Usage:
  botio server [flags]

Examples:
botio server --db ./data/botio.db --col commands --http :9090 --key mysupersecretkey

Flags:
      --col string       collection used to store the commands (default "commands")
      --db string        path to the database (default "./botio/botio.db")
  -h, --help             help for server
      --http string      port for HTTP connections (default ":80")
      --https string     port for HTTPS connections (default ":443")
      --key string       authentication key for JWT
      --sslcert string   ssl certification
      --sslkey string    ssl key
```

Example:

```bash
botio server --db ./data/botio.db --col commands --http :9090 --key mysupersecretkey
```

> The database used is based on BoltDB. You can read more about it [here](https://github.com/etcd-io/bbolt).

### telegram

```text
$ botio telegram --help
Initializes a Telegram bot that extracts the commands from the botio's server.

Usage:
  botio telegram [flags]

Examples:
botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey

Flags:
  -h, --help           help for telegram
      --key string     authentication key for JWT
      --token string   telegram's token
      --url string     botio's server URL
```

Example:

```bash
botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey
```

### discord

```text
$ botio discord --help
Initializes a Discord bot that extracts the commands from the botio's server

Usage:
  botio discord [flags]

Examples:
botio discord --token <discord-token> --url :9090 --key mysupersecretkey

Flags:
  -h, --help           help for discord
      --key string     authentication key for JWT
      --token string   discord's token
      --url string     botio's server URL
```

Example:

```bash
botio discord --token <discord-token> --url :9090 --key mysupersecretkey
```

### add

```text
$ botio add --help
Adds a new command with a response to the botio's server

Usage:
  botio add [flags]

Examples:
botio add --command start --response Hello --url :9090 --password mypassword

Flags:
      --command string    command to add
  -h, --help              help for add
      --key string        authentication key for JWT
      --response string   response of the command to add
      --url string        url where the botio's server is listening
```

Example:

```bash
bodio add --command start --response Hello --url :9090 --key mysupersecretpassword
```

### print

```text
$ botio print
Prints the specified botio's command with his response

Usage:
  botio print [flags]

Examples:
botio print --command start --url :9090 --key mysupersecretkey

Flags:
      --command string   command to search for
  -h, --help             help for print
      --key string       authentication key for JWT
      --url string       url where the botio's server is listening
```

Example:

```bash
botio print --command start --url :9090 --key mysupersecretkey
```

### list

```text
$ botio list --help
Prints a list with all the botio's commands

Usage:
  botio list [flags]

Examples:
botio list --url :9090 --key mysupersecretkey

Flags:
  -h, --help         help for list
      --key string   authentication key for JWT
      --url string   url where the botio's server is listening
```

Example:

```bash
botio list --url :9090 --key mysupersecretkey
```

### update

```text
$ botio update --help
Updates an existing command (or adds it if not exists) with a response on the botio's server

Usage:
  botio update [flags]

Examples:
botio update --command start --response Hi --url :9090 --key mysupersecretkey

Flags:
      --command string    command to add
  -h, --help              help for update
      --key string        authentication key for JWT
      --response string   response of the command to add
      --url string        url where the botio's server is listening
```

Example:

```text
botio update --command start --response Hi --url :9090 --key mysupersecretkey
```

### delete

```text
$ botio delete --help
Deletes the specified botio's command from the botio's server

Usage:
  botio delete [flags]

Examples:
botio delete --command start --url :9090 --key mysupersecretkey

Flags:
      --command string   command to delete
  -h, --help             help for delete
      --key string       authentication key for JWT
      --url string       url where the botio's server is listening
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
