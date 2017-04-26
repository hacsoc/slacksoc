# slacksoc

Our friendly little slackbot.

## Features

* Is very friendly (It will say hi back to you)
* Sends you a private message on request
* Welcomes new members by requesting that they set their Real Name field

## Development

Nearly all development of bots and plugins takes place in the [library][lib].
This repository merely contains deployment information (e.g. Dockerfile),
configuration, and pinned dependencies.

To use the pinned dependencies, be sure to `git submodule init --update`
initially, and `git submodule update` each time you update the repository. Go
will use code within the `vendor` directory rather than downloading the latest
versions. `go build` will build the bot, and `docker build` will create a Docker
image.

## Updating Dependencies

The convenient way to update all submodules is to use [vendetta][]. Install it
with `go get vendetta`, and then use `vendetta -u` to update the submodules. You
can also use `vendetta -p` to prune unneeded dependencies. Then you can commit
the updated submodules and build a new Docker image.

[lib]: https://github.com/brenns10/slacksoc
[vendetta]: https://github.com/dwp/vendetta
