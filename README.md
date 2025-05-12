# aagh -- ![GitHub License](https://img.shields.io/github/license/kermage/aagh) ![GitHub Release](https://img.shields.io/github/v/release/kermage/aagh) ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/kermage/aagh/total)

> To the betterment of a project: no more **"team, please lint and test"** and **"it works on my machine"**. _Hopefully_.

```
A cross-platform executable for handling Git hooks

Usage:
  aagh [command]

Available Commands:
  check       Check the repository status in the current directory
  init        Initialize the repository in the current directory
  run         Run a hook in the repository of current directory
  setup       Setup a hook in the repository of current directory
  help        Help about any command
  completion  Generate the autocompletion script for the specified shell

Flags:
  -h, --help      help for aagh
  -v, --version   version for aagh

Use "aagh [command] --help" for more information about a command.
```

## Installation

### Prepared Binaries

Download the latest release [here](https://github.com/kermage/aagh/releases).

### One-liner Commands

```sh
curl -fsSL https://raw.githubusercontent.com/kermage/aagh/main/install.sh | sh
```

```sh
wget -qO- https://raw.githubusercontent.com/kermage/aagh/main/install.sh | sh
```

*Default install path: `/usr/local/bin`*

#### Custom Bin Directory

```sh
sh -c "$(curl -fsSL https://raw.githubusercontent.com/kermage/aagh/main/install.sh)" -- -b <path>
```

```sh
sh -c "$(wget -qO- https://raw.githubusercontent.com/kermage/aagh/main/install.sh)" -- -b <path>
```

*\* Make sure `<path>` exists and is writeable. Directory creation is intentionally not implemented.*

### Project Examples

- [GO](./examples/go)
- [JS](./examples/js)
- [PHP](./examples/php)

#### Oneshot command

Simplify project hooks setup

```sh
aagh init --apply
```
