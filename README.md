# Gopher Street

<image align="right" width="200" height="200" src="gopher.svg" title="Credit: gopherize.me">

![workflow](https://github.com/amitlevy21/gopher-street/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/amitlevy21/gopher-street)](https://goreportcard.com/report/github.com/amitlevy21/gopher-street)
[![codecov](https://codecov.io/gh/amitlevy21/gopher-street/branch/main/graph/badge.svg?token=y0n7I2Ind3)](https://codecov.io/gh/amitlevy21/gopher-street)

Do you know how much money you spend per month? On what?

What are the costs of such expenses on the long term?

How can we monitor and understand where our money is going?

Gopher Street is here to help.

## Features

- Reads from your bank exported files, instead of manually have to input them one by one
- Reads transaction for multiple sources, storing them in a single centralized place
- Connects to MongoDB database, which can be self managed by the user
- Supports `.csv` and `.xlsx` formats
- More to come!

## Installation

```sh
go get github.com/amitlevy21/gopher-street
```

## Usage

```sh
❯ gst --help

Usage:
  gst [command]

Available Commands:
  get         Get expenses from DB
  help        Help about any command
  load        Load data from file

Flags:
  -h, --help   help for gst

Use "gst [command] --help" for more information about a command.
```

## Development

### Configure Git Hooks

Install [pre-commit](https://pre-commit.com) and [commitizen](https://commitizen-tools.github.io/commitizen/).

```sh
❯ pre-commit install
❯ pre-commit install --hook-type commit-msg
❯ pre-commit install-hooks
```

Commits messages are formatted with `commitizen`:

```sh
❯ cz commit
```

## License

[MIT License](./License)
