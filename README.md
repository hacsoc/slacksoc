# slacksoc

[![Build Status](https://travis-ci.org/hacsoc/slacksoc.svg?branch=master)](https://travis-ci.org/hacsoc/slacksoc)
[![Coverage Status](https://coveralls.io/repos/github/hacsoc/slacksoc/badge.svg?branch=coverage)](https://coveralls.io/github/hacsoc/slacksoc?branch=coverage)

Our friendly little slackbot.

## Features

* Is very friendly (It will say hi back to you)
* Sends you a private message on request
* Welcomes new members by requesting that they set their Real Name field
* Does not reply to other bots, avoiding botpocalypse

## Development

You'll need a go environment set up on your machine, which you can get
[here](https://golang.org/). Just click the big "Download Go" button and follow
the instructions.

If you want to run the actual bot locally, you'd need the bot's API token,
which I probably won't give you. Luckily, you can write tests! `go test` will
run the test suite.
