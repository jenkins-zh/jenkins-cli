[简体中文](https://github.com/jenkins-zh/jenkins-cli/blob/master/README-zh.md)

# Jenkins CLI

[![Go Report Card][go-report-card-badge]][go-report-card-url]
[![Quality Gate Status][sonar-badge]][sonar-link]
[![GoDoc][godoc-badge]][godoc-url]
![Sonar Coverage](https://img.shields.io/sonar/coverage/jenkins-zh_jenkins-cli?server=https%3A%2F%2Fsonarcloud.io)
[![Travis](https://img.shields.io/travis/jenkins-zh/jenkins-cli.svg?logo=travis&label=build&logoColor=white)](https://travis-ci.org/jenkins-zh/jenkins-cli)
[![Contributors](https://img.shields.io/github/contributors/jenkins-zh/jenkins-cli.svg)](https://github.com/jenkins-zh/jenkins-cli/graphs/contributors)
[![GitHub release](https://img.shields.io/github/release/jenkins-zh/jenkins-cli.svg?label=release)](https://github.com/jenkins-zh/jenkins-cli/releases/latest)
![GitHub All Releases](https://img.shields.io/github/downloads/jenkins-zh/jenkins-cli/total)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/jenkins-zh/jenkins-cli)
[![Gitter](https://badges.gitter.im/jenkinsci/jenkins-cli.svg)](https://gitter.im/jenkinsci/jenkins-cli?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

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

Read [this document](doc/README.md) to know more details about how to use `jcli`.

# Contribution

If you're interested in this project. Please go through the
[contribution guide](CONTRIBUTING.md). Any contributions are welcome.

# Stargazers over time

[![Stargazers over time](https://starchart.cc/jenkins-zh/jenkins-cli.svg)](https://starchart.cc/jenkins-zh/jenkins-cli)

[go-report-card-url]: https://goreportcard.com/report/jenkins-zh/jenkins-cli
[go-report-card-badge]: https://goreportcard.com/badge/jenkins-zh/jenkins-cli
[sonar-badge]: https://sonarcloud.io/api/project_badges/measure?project=jenkins-zh_jenkins-cli&metric=alert_status
[sonar-link]: https://sonarcloud.io/dashboard?id=jenkins-zh_jenkins-cli
[godoc-url]: https://godoc.org/github.com/jenkins-zh/jenkins-cli
[godoc-badge]: http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square
