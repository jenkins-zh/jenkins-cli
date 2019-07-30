# Jenkins CLI

[![Go Report Card](https://goreportcard.com/badge/jenkins-zh/jenkins-cli)](https://goreportcard.com/report/jenkins-zh/jenkins-cli)

Jenkins CLI allows you manage your Jenkins as an easy way. No matter you're a plugin
developer, administrator or just a regular user, it borns for you!

# Features

* Multiple Jenkins support
* Plugins management (list, search, install, upload)
* Job management (search, build, log)
* Open your Jenkins with a browse
* Restart your Jenkins
* Connection with proxy support

# Get it

We support mac, linux and windows for now.

## mac

You can use `brew` to install jcli.
```
brew tap jenkins-zh/jcli
brew install jcli
```

## Linux

It's very simple to install `jcli` into your Linux OS. Just need to execute a command line at below:
```
curl -L https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-linux-amd64.tar.gz|tar xzv
sudo mv jcli /usr/local/bin/
```

## Windows

You can find the latest version by [click here](https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-windows-386.tar.gz). Then download the tar file, cp the uncompressed `jcli` directory into your system path.

## Other Package Managers

Here're other package managers:

* [GoFish](https://gofi.sh/) users can use `gofish install jcli`.

# Get started

Once you'v installed `jcli`. You should provide a config for it. Please execute cmd `jcli config generate`, this will help you to edit the config file `~/.jenkins-cli.yaml`. According to your Jenkins configuration to modify this file.

# Contribution

It's still under very early develope time. Any contribution is welcome.
