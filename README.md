# Weave
Write detailed commit messages without leaving the terminal

## Why
- You need to write a description for your commit
- You want to do it from the terminal

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Configuation](#configuration)

## Installation
### Go
```sh
go install github.com/joeldotdias/weave@latest
```

### From Source
```
git clone https://github.com/joeldotdias/weave
cd weave
go install
```

## Usage
```sh
weave
```

### Flags
How your title is processed<br>
```
fix (🛠️): Fixed bug
<type> (<symbol>)<separator>Title...
```
This can be changed in the [configuration](#configuration) file or via flags
`-a, --add_all` Adds all files that have not yet been staged to commit<br>
`-t, --title` Lets you provide a title and directly skip to the description<br>
`-f, --format` Format for the type and symbol in the prefix<br>
`-s, --separator` The separator between the prefix and title. The separator is placed without any additional formatting so whitespace is significant

## Configuration
Weave looks for the config in
```
$XDG_CONFIG_HOME/weave/config.toml
or
$HOME/.config/weave/config.toml
```
The prefixes along with the symbols are stored here. Only the prefixes provided by you will be used.
Weave also supports per project configuration. For this a file called `weave.toml` must be placed in the root directory of the project

### Example
#### config.toml
```toml
add_all = false
format = "<type> <symbol>"
separator = ": "

[symbols]
feat = "🚀"
fix = "🛠️"
chore = "📦"
```
