# gclone

Sick of having unorganised repositories? Tired of specifying the directory?
Get `gclone` today!

```
/home/samwise/code
├── bitbucket.org
│   └── company
│       ├── alarms
│       ├── dashboards
│       ├── product-alpha
│       ├── product-bravo
│       ├── references
│       ├── tooling
│       └── utility-belt
├── github.com
│   └── scottgreenup
│       ├── gclone
│       ├── json-util
│       └── lolcat
└── gitlab.com
    ├── graphviz
    │   └── graphviz
    └── wireshark
        └── wireshark
```

## Installation

To get the latest version:

```shell
go install github.com/scottgreenup/gclone@latest
```

## Example usage

```text
[~]# gclone git@github.com:scottgreenup/desktop.git
git clone git@github.com:scottgreenup/desktop.git /Users/samwise/code/github.com/scottgreenup/desktop
Cloning into '/Users/samwise/code/github.com/scottgreenup/desktop'...
remote: Enumerating objects: 839, done.
remote: Counting objects: 100% (112/112), done.
remote: Compressing objects: 100% (14/14), done.
remote: Total 839 (delta 100), reused 98 (delta 98), pack-reused 727 (from 1)
Receiving objects: 100% (839/839), 10.74 MiB | 4.99 MiB/s, done.
Resolving deltas: 100% (347/347), done.
{"targetDirectory":"/Users/samwise/code/github.com/scottgreenup/desktop"}

[~]# cd /Users/samwise/code/github.com/scottgreenup/desktop

[~/code/github.com/scottgreenup/desktop]# git status
On branch master
Your branch is up to date with 'origin/master'.

nothing to commit, working tree clean
```

You can use it with regular git arguments if you need:

```shell
# gclone https://github.com/scottgreenup/desktop.git ./here -- --no-checkout
git clone --no-checkout https://github.com/kubernetes/kubernetes.git here
...
```

## Automatically change directory

### Script

You can also create a script so you can enter the directory immediately.
Ensure you have your `$PATH` setup to prioritise the bash script, then:

```shell
$ . gclone https://github.com/scottgreenup/desktop.git
$ pwd
/home/scottgreenup/code/github.com/scottgreenup/desktop
```

### ZSH

You can slap this in your `.zshrc`

```shell
_gclone_binary=$(command -v gclone)
function gclone {
    local _repo="${1}"
    cd $( ${_gclone_binary} "${_repo}" | jq -r .targetDirectory)
}
```

This one is more robust:

```shell
_gclone_binary=$(command -v gclone)
if [[ ! -z "${_gclone_binary}" ]]; then
    function gclone {
        local _repo="${1}"
        cd $( ${_gclone_binary} "${_repo}" | jq -r .targetDirectory)
    }
else
    1>&2 echo "unable to find gclone binary"
fi
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

