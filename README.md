# osx-remap-keyboard-modifiers

This is a tool to set keyboard modifiers on OSX. Intended for use in "sane
defaults" scripts and what-not to disable the horrific caps lock key and
repurpose it as `INSERT YOUR FAVORITE MOD KEY HERE`

## Install

Assuming you have a Go development environment set up,
```shell
go get github.com/rtlong/osx-remap-keyboard-modifiers
```

If I get this cleaned up and tested properly, I'll try to get it in the main Homebrew repo for easier install.

## Usage

For example, to remap the Caps Lock key to act as though it were the left Control key, do this:
```shell
osx-remap-keyboard-modifiers caps:control_l
```

## Contributions

PRs and Issues are welcome!
