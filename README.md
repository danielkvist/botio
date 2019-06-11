# Botio

[![Go Report Card](https://goreportcard.com/badge/github.com/danielkvist/botio)](https://goreportcard.com/report/github.com/danielkvist/botio)
[![Docker Pulls](https://img.shields.io/docker/pulls/danielkvist/botio.svg?maxAge=604800)](https://hub.docker.com/r/danielkvist/botio/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

Botio is a quick and easy way to create a simple and easy to update and maintain Telegram bot.

## Example

You can see an example [here](https://t.me/dkvist_bot).

> P.S.: I'm looking for a job as a Junior Go Developer.

## How it works?

Botio basically creates a BoltDB database if it does not exist in which all the commands and their responses for the Telegram bot will be stored. This database can be updated through a simple CRUD API that starts alongside the Telegram bot.

## ToDo

- [] Better logging.
- [] Web Interface.

## Installation

I recommend to use Docker:

```bash
docker image pull danielkvist/botio
```

But you can also use ```go install```:

```bash
go get github.com/danielkvist/botio
```

Or clone my GitHub repository:

```bash
git clone https://github.com/danielkvist/botio.git botio
```

## Configuration

To set up Botio and get started you just need to create an ```.env``` file.

```text
TELEGRAM_TOKEN=<your-telegram-token>
DATABASE=<path-of/for-the-boltdb-database>
COLLECTION=<name-of-the-collection-for-the-commands>
LISTEN_ADDRESS=<address-for-your-api>
API_USERNAME=<username-for-your-api>
API_PASSWORD=<password-for-your-api>
```

For example:

```text
TELEGRAM_TOKEN=42...
DATABASE=/data/botio.db
COLLECTION=commands
LISTEN_ADDRESS=:9090
API_USERNAME=user
API_PASSWORD=password
```

## Start

```bash
docker container run --name my-bot --env-file .env -v /data:/data:rw -p 9090:9090 danielkvist/botio
```

## CRUD API

As I said before Botio provides a simple CRUD API to manage the database from which the bot for Telegram extracts the commands and their corresponding responses.

> The following examples will use ```curl```. But you can also use Postman, for example.

### GET

```bash
curl -u user:password -X GET http://localhost:9090/api/commands/start
# Response: {"cmd":"start","response":"Hi!"}
```

### GET All

```bash
curl -u user:password -X GET http://localhost:9090/api/commands
# Response: [{"cmd":"start","response":"Hi!"},{"cmd":"goodbye","response":"I see you later!"},...]
```

### POST

```bash
echo '{"cmd": "age", "response":"42"}' | curl -u user:password -d @- http://localhost:9090/api/commands
# Response: {"cmd":"age","response":"42"}
```

### PUT

> In a Bolt database updating is the same as reposting the same element but with a different value. 

```bash
echo '{"cmd": "age", "response":"25"}' | curl -u user:password -d @- http://localhost:9090/api/commands
# Response: {"cmd": "age","response":"25"}
```

### DELETE

```bash
curl -u user:password -X DELETE http://localhost:9090/api/commands/start
```

### Backup

```bash
curl -u user:password -X GET http://localhost:9090/api/backup > backup.db
```