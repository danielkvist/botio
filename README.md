# Botio

[![Go Report Card](https://goreportcard.com/badge/github.com/danielkvist/botio)](https://goreportcard.com/report/github.com/danielkvist/botio)
[![Docker Pulls](https://img.shields.io/docker/pulls/danielkvist/botio.svg?maxAge=604800)](https://hub.docker.com/r/danielkvist/botio/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

Botio is a CLI to create and manage easily simple bots for differents platforms.

> For the moment it only supports Telegram but I'm working to add support for more platforms soon.

## Example

First, we create a server from which we're going to manage the commands to our bot:

```bash
botio server --db ./data/botio.db --col commands --addr localhost:9090 --user myuser --password mypassword
```

We add a command:

```bash
botio add --command start --response Hello --url localhost:9090 --user myuser --password mypassword
```

We can check that the command has been created successfully with ```print```:

```bash
botio print --command start --url localhost:9090 --user myuser --password mypassword
```

Or see a list of all our commands with ```list```:

```bash
botio list --url localhost:9090 --user myuser --password mypassword
```

Now, we can start our Telegram's bot:

```bash
botio tbot --token <telegram-token> --url localhost:9090 --user myuser --password mypassword
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
botio server --db ./data/commands.db --col commands --addr localhost:9090 --user myuser --password mypassword
botio tbot --token <telegram-token> --url localhost:9090 --user myuser --password mypassword
botio print --command start --url localhost:9090 --user myuser --password mypassword

Available Commands:
  add         Adds a new command with a response to the botio's server
  delete      Deletes the specified botio's command from the botio's server
  help        Help about any command
  list        Prints a list with all the botio's commands
  print       Prints the specified botio's command with his response
  server      Starts a botio's server to manage the botio's commands with simple HTTP methods.
  tbot        Initializes a Telegram's bot that extracts the commands from the botio's server.
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
botio server --db ./data/botio.db --col commands --addr localhost:9090 --user mysuer --password mypassword

Flags:
      --addr string       address where the server should listen for requests (default "localhost:9090")
      --col string        collection used to store the commands (default "commands")
      --db string         path to the database (default "./botio/botio.db")
  -h, --help              help for server
      --password string   password for basic auth (default "toor")
      --user string       username for basic auth (default "admin")
```

Example:

```bash
botio server --db ./data/botio.db --col commands --addr localhost:9090 --user myuser --password mypassword
```

> The database used is based on BoltDB. You can read more about it [here](https://github.com/etcd-io/bbolt).

### tbot

```text
$ botio tbot --help
Initializes a Telegram's bot that extracts the commands from the botio's server.

Usage:
  botio tbot [flags]

Examples:
botio tbot --token <telegram-token> --url localhost:9090 --user myuser --password mypassword

Flags:
  -h, --help              help for tbot
      --password string   password for basic auth (default "toor")
      --token string      Telegram's token
      --url string        URL where the botio's server is listening for requests
      --user string       username for basic auth (default "admin")
```

Example:

```bash
botio tbot --token <telegram-token> --url localhost:9090 --user myuser --password mypassword
```

### add

```text
$ botio add --help
Adds a new command with a response to the botio's server

Usage:
  botio add [flags]

Examples:
botio add --command start --response Hello --url localhost:9090 --user myuser --password mypassword

Flags:
      --command string    command to add
  -h, --help              help for add
      --password string   password for basic auth (default "toor")
      --response string   response of the command to add
      --url string        URL where the botio's server is listening
      --user string       username for basic auth (default "admin")
```

Example:

```bash
bodio add --command start --response Hello --url localhost:9090 --user myuser --password mypassword
```

### print

```text
$ botio print
Prints the specified botio's command with his response

Usage:
  botio print [flags]

Examples:
botio print --command start --url localhost:9090 --user myuser --password mypassword

Flags:
      --command string    command to search for
  -h, --help              help for print
      --password string   password for basic auth (default "toor")
      --url string        URL where the botio's server is listening
      --user string       username for basic auth (default "admin")
```

Example:

```bash
botio print --command start --url localhost:9090 --user myuser --password mypassword
```

### list

```text
$ botio list --help
Prints a list with all the botio's commands

Usage:
  botio list [flags]

Examples:
botio list --url localhost:9090 --user myuser --password mypassword

Flags:
  -h, --help              help for list
      --password string   password for basic auth (default "toor")
      --url string        URL where the botio's server is listening
      --user string       username for basic auth (default "admin")
```

Example:

```bash
botio list --url localhost:9090 --user myuser --password mypassword
```

### update

```text
$ botio update --help
Updates an existing command (or adds it if not exists) with a response on the botio's server

Usage:
  botio update [flags]

Examples:
botio update --command start --response Hi --url localhost:9090 --user myuser --password mypassword

Flags:
      --command string    command to add
  -h, --help              help for update
      --password string   password for basic auth (default "toor")
      --response string   response of the command to add
      --url string        URL where the botio's server is listening
      --user string       username for basic auth (default "admin")
```

Example:

```text
botio update --command start --response Hi --url localhost:9090 --user myuser --password mypassword
```

### delete

```text
$ botio delete --help
Deletes the specified botio's command from the botio's server

Usage:
  botio delete [flags]

Examples:
botio delete --command start --url localhost:9090 --user myuser --password mypassword

Flags:
      --command string    command to delete
  -h, --help              help for delete
      --password string   password for basic auth (default "toor")
      --url string        URL where the botio's server is listening
      --user string       username for basic auth (default "admin")
```

Example:

```bash
botio delete --command start --url localhost:9090 --user myuser --password mypassword
```

## API Endpoints

### GET

#### GET Commands

```text
http://<url>:<port>/api/commands
```

Example:

```text
http://localhost:9090/api/commands
```

```bash
curl -u myuser:mypassword -X GET http://localhost:9090/api/commands
```

#### GET Command

```text
http://<url>:<port>/api/commands/<command>
```

Example:

```text
http://localhost:9090/api/commands/start
```

```bash
curl -u myuser:mypassword -X GET http://localhost:9090/api/commands/start
```

### POST

```text
http://<url>:<port>/api/commands
```

Example:

```text
http://localhost:9090/api/commands
```

```bash
echo '{"cmd": "age", "response":"42"}' | curl -u myuser:mypassword -d @- http://localhost:9090/api/commands
```

### UPDATE

```text
http://<url>:<port>/api/commands/<command>
```

Example:

```text
http://localhost:9090/api/commands/start
```

```bash
echo '{"cmd": "age", "response":"25"}' | curl -u myuser:mypassword -d @- http://localhost:9090/api/commands/start
```

### DELETE

```text
http://<url>:<port>/api/commands/<command>
```

Example:

```text
http://localhost:9090/api/commands/start
```

```bash
curl -u myuser:mypassword -X DELETE http://localhost:9090/api/commands/start
```

### Backup

> Backup is an special endpoint that will send to the client a backup of the database.

```text
http://<url>:<port>/api/backup
```

Example:

```text
http://localhost:9090/api/backup
```

```bash
curl -u myuser:mypassword -X GET http://localhost:9090/api/backup > backup.db
```

## ToDo

- [ ] Better logging
- [ ] Server with HTTPS
- [ ] Docker Compose
- [ ] Web Interface
- [ ] Alternative databases like PostgreSQL
- [ ] Alternative authentication options
- [ ] Support for Facebook Messenger bots
- [ ] Support for Discord bots
- [ ] Support for Slack bots
- [ ] Support for Skype bots
