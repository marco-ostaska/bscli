# bscli

bscli is a command line interface to interact with [bluesight](https://www.bluesight.io/) kanban.

[![GoDoc](https://godoc.org/github.com/marco-ostaska/bscli?status.svg)](https://godoc.org/github.com/marco-ostaska/bscli)
[![Go Report Card](https://goreportcard.com/badge/github.com/marco-ostaska/bscli)](https://goreportcard.com/report/github.com/marco-ostaska/bscli)

# Table of Contents

- [Overview](#overview)
- [Commands](#commands)
- [Installing](#intalling)
- [Getting Started](#getting-started)


# Overview

bscli is a command line interface provides an interface to interact over daily tasks administration using [bluesight](https://www.bluesight.io/) kanban.

# Commands

- [bscli](docs/bscli.md) - A command line tool for bluesight.io
  - [vault](docs/bscli_vault.md) - create or update vault credentials.
  - [squads](docs/bscli_squads.md) - list the squads for the user
  - [squad](docs/bscli_squad.md) - display information for a given squad
  - [mood](docs/bscli_mood.md) - display mood marbles information for a given squad
  - [card](docs/bscli_card.md) - create or update cards
  

# Getting Started

Before using the bscli to perform admnistrative tasks you must [configure the bscli vault](docs/cmd/bscli_vault_new.md).

For more information how to create a token, [please read this first](https://portal.bluesight.io/tutorial.html) under API section.

# License

bscli is released under the GNU General Public License v3. See [LICENSE](LICENSE)
