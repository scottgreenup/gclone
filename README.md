# gclone

Sick of having unorganised repositories? Tired of specifying the directory?
Get `gclone` today!

## Installation

To get the latest version:

```shell
go install github.com/scottgreenup/gclone@latest
```

## Example usage

```shell
$ gclone https://github.com/scottgreenup/desktop.git
$ cd ~/code/github.com/scottgreenup/desktop
$ pwd
/home/scottgreenup/code/github.com/scottgreenup/desktop

$ git status
On branch master
Your branch is up to date with 'origin/master'.

nothing to commit, working tree clean
```

You can use it with args:

```shell
$ gclone https://github.com/scottgreenup/desktop.git ./here -- --no-checkout
git clone --no-checkout https://github.com/kubernetes/kubernetes.git here
...
```

## Automatically change directory

You can also create a script so you can enter the directory immediately.

```shell
$ cat ~/bin/gclone 
#!/usr/bin/env bash

cloned_directory=$(~/go/bin/gclone $@ | jq -r .targetDirectory)
cd $cloned_directory
```

Ensure you have your `$PATH` setup to prioritise the bash script, then:

```shell
$ . gclone https://github.com/scottgreenup/desktop.git
$ pwd
/home/scottgreenup/code/github.com/scottgreenup/desktop
```

## Configuration

You can configure `gclone` via a configuration file.

* `$HOME/.config/gclone/config.json`
* `$HOME/.config/gclone/config.yaml`
* `/etc/gclone/config.json`
* `/etc/gclone/config.yaml`

### Configuration options

#### DefaultDirectory

Default value is `~/code`

The default directory to clone into. Ensure it is created before using as
`glone` will not create it for you.

### Example configuration

```json
{
  "DefaultDirectory": "~/dev/"
}
```

