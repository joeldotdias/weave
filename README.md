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
`-a, --add_all` Adds all files that have not yet been staged to commit
`-t, --title` Lets you provide a title and directly skip to the description
