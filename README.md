# glcone

Sick of having unorganised repositories? Tired of specifying the directory?
Get `gclone` today!

## Example Usage

```
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

```
$ gclone https://github.com/scottgreenup/desktop.git ./here -- --no-checkout
git clone --no-checkout https://github.com/kubernetes/kubernetes.git here
...
```
