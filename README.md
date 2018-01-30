Coaching - Slack bot
====================

# Pre-requisites

You need to have `dep` tool installed to manage dependencies: [https://github.com/golang/dep](https://github.com/golang/dep).

# Installation

First, you need to install dependencies by running:

```
$ make install
```

Then, you have to create a `.env` file with your own values (you can find these values on Slack under #coach's channel topic):

```
$ cp .env.dist .env
<edit .env file to fill your Slack token and google spreadsheet id>
$ source .env
```

You are almost ready to go! You need a last step: create a Google API OAuth 2.0 token:

You have to follow the instructions located in the link below in order to generate OAuth 2.0 Client ID keys: [https://cloud.google.com/genomics/downloading-credentials-for-api-access](https://cloud.google.com/genomics/downloading-credentials-for-api-access)

Once it's done, download the `client_secret.json` file from the interface and put it under `config/` folder.

On the first launch of the bot, you will also have to authenticate using OAuth2 with Google APIs. Just click on the link it outputs.

# Usage (on dev environment)

```
$ go run main.go
```

# Compilation

Run `make build` to build the source into a single `elevenbot` binary file.

# Prepare deployment on a new server

On first time deployment, you have to:

Edit your `.env` file in order to edit `ELEVENBOT_SSH_USER` and `ELEVENBOT_SSH_IP` environment variables.
Edit `config/elevenbot.service` file with environment variables and user name.

Then, run `make prepare-deploy` to prepare your instance.

# Deploy

Just run `make deploy`