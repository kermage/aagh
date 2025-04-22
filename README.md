# aagh --

> To the betterment of a project: no more **"team, please lint and test"** and **"it works on my machine"**. _Hopefully_.

```
A cross-platform executable for handling Git hooks

Usage:
  aagh [command]

Available Commands:
  check       Check the repository status in the current directory
  init        Initialize the repository in the current directory
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

### Project Examples

- [GO](./examples/go)
- [JS](./examples/js)
- [PHP](./examples/php)

#### Oneshot command

Simplify project hooks setup

```sh
aagh init --apply
```
