# bscli

bscli is a command line interface to interact with [bluesight](https://www.bluesight.io/) kanban.

[![Go Report Card](https://goreportcard.com/badge/github.com/marco-ostaska/bscli)](https://goreportcard.com/report/github.com/marco-ostaska/bscli)
[![Build Status](https://travis-ci.com/marco-ostaska/bscli.svg?branch=unreleased)](https://travis-ci.com/marco-ostaska/bscli)

# Table of Contents

- [Overview](#overview)
- [Commands](#commands)
- [Installing](#intalling)
- [Getting Started](#getting-started)
- [Limitations](#limitations)
- [License](#license)


# Overview

bscli is a command line interface provides an interface to interact over daily tasks administration using [bluesight](https://www.bluesight.io/) kanban.

# Commands

- [bscli](docs/bscli.md) - A command line tool for bluesight.io
  - [vault](docs/bscli_vault.md) - create or update vault credentials.
  - [squads](docs/bscli_squads.md) - list the squads for the user
  - [squad](docs/bscli_squad.md) - display information for a given squad
  - [card](docs/bscli_card.md) - create or update cards

# Install

[download](https://github.com/marco-ostaska/bscli/releases) the binary

Put it in a safe place and make it executable

```shell
# For linux and mac only
chmod +x bscli-<platform>-<version>
```

Then simply execute the app.

# Getting Started

Before using the bscli to perform admnistrative tasks you must [configure the bscli vault](docs/bscli_vault_new.md)).

For more information how to create a token, [please read this first](https://portal.bluesight.io/tutorial.html) under API section.

# Limitations

- bscli can not overlap bluesight.io api limitations such as adding mood with it.
- bscli does not provide all functionalities that bluesight.io does, with provide basic functionality for daily uses, even if it is provided by the API.
- Feel free to raise a [feature request](https://github.com/marco-ostaska/bscli/issues/new?assignees=&labels=&template=feature_request.md&title=) if you believe something is missing and if it is not an API limitation we will be glad to verify the implementation possibility.

# License

bscli is released under the GNU General Public License v3. See [LICENSE](LICENSE)
