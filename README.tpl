# Quick start

[简体中文](https://github.com/jenkins-zh/jenkins-cli/blob/master/README-zh.md)

## Jenkins CLI

<!--
[![](https://sonarcloud.io/api/project_badges/measure?project=jenkins-zh_jenkins-cli&metric=alert_status)](https://sonarcloud.io/dashboard?id=jenkins-zh_jenkins-cli) 
-->
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/jenkins-zh/jenkins-cli)
[![](https://goreportcard.com/badge/jenkins-zh/jenkins-cli)](https://goreportcard.com/report/jenkins-zh/jenkins-cli)
[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/jenkins-zh/jenkins-cli)
[![codecov](https://codecov.io/gh/jenkins-zh/jenkins-cli/branch/master/graph/badge.svg?token=XS8g2CjdNL)](https://codecov.io/gh/jenkins-zh/jenkins-cli)
[![Contributors](https://img.shields.io/github/contributors/jenkins-zh/jenkins-cli.svg)](https://github.com/jenkins-zh/jenkins-cli/graphs/contributors)
[![GitHub release](https://img.shields.io/github/release/jenkins-zh/jenkins-cli.svg?label=release)](https://github.com/jenkins-zh/jenkins-cli/releases/latest)
![GitHub All Releases](https://img.shields.io/github/downloads/jenkins-zh/jenkins-cli/total)
[![Docker Pulls](https://img.shields.io/docker/pulls/jenkinszh/jcli.svg)](https://hub.docker.com/r/jenkinszh/jcli/tags)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/jenkins-zh/jenkins-cli)
[![HitCount](http://hits.dwyl.com/jenkins-zh/jenkins-cli.svg)](http://hits.dwyl.com/jenkins-zh/jenkins-cli)

Jenkins CLI allows you manage your Jenkins in an easy way. No matter if you're a plugin developer, administrator or just a regular user, it is made for you!

{{printToc}}

## Features

* Multiple Jenkins support
* Plugins management \(list, search, install, upload\)
* Job management \(search, build, log\)
* Configuration as Code support
* Open your Jenkins with a browser
* Restart your Jenkins
* Connection with proxy support

## Get it

We support macOS, Linux and Windows for now.
<a href="https://repology.org/project/jenkins-cli/versions">
    <img src="https://repology.org/badge/vertical-allrepos/jenkins-cli.svg" alt="Packaging status" align="right">
</a>

### macOS

You can use `brew` to install jcli.

```text
brew tap jenkins-zh/jcli
brew install jcli
```

Alternatively, you can use [MacPorts](https://ports.macports.org/port/jenkins-cli/summary).

```text
sudo port install jenkins-cli
```

For MacPorts, add `+bash_completion` or `+zsh_completion` for shell completion.

### Linux

To install `jcli` on your Linux OS, execute the following command:

```sh
curl -L https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-linux-amd64.tar.gz|tar xzv
sudo mv jcli /usr/local/bin/
```

### Windows

You can install `jcli` via [scoop](https://scoop.sh/) or [choco](https://chocolatey.org/packages/jcli/). 

Or you can also find the latest version from the [release page](https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-windows-amd64.zip). 
Download the zip file and copy the uncompressed `jcli` directory into your system path.

### Other package managers

Here are other package managers:

| Install | Upgrade | Uninstall | Platform |
|---|---|---|---|
| `scoop install jcli` | | | `Windows` |
| `choco install jcli` | `choco upgrade jcli` | `choco uninstall jcli` | `Windows` |
| `snap install jcli` | `snap refresh jcli` | `snap remove jcli` | `Linux` |
| `sudo port install jenkins-cli`| `sudo port selfupdate && sudo port upgrade jenkins-cli` | `sudo port uninstall jenkins-cli` | `macOS`

See more about [how to download jcli](docs/book/en/download.md). You can find the download details [from here](https://tooomm.github.io/github-release-stats/?username=jenkins-zh&repository=jenkins-cli).

## Get started

Read the [official document](http://jcli.jenkins-zh.cn/) for more details on how to use `jcli`.

Or, you can take [a live interactive course](https://www.katacoda.com/jenkins-zh/scenarios/course-jcli) of Jenkins CLI.

## Plugins

Jenkins CLI could have more features by installing a plugin for it. You can install a plugin by the following command:

```text
jcli config plugin fetch
jcli config plugin install account
```

All official plugins could be found at [here](https://github.com/jenkins-zh/jcli-plugins).

## Contribution

If you're interested in this project. Please go through the [contribution guide](https://github.com/jenkins-zh/jenkins-cli/tree/cb3d358df4699db11b681eb0ab9adffbfb8a7bd4/CONTRIBUTING.md). Any contributions are welcome.

Thanks to JetBrains for giving us the open source licence.  
[![goland.svg](docs/book/.gitbook/assets/goland.svg)](https://www.jetbrains.com/?from=jenkins-cli)

## Similar Projects

There're a few similar projects that you might be interested in:

* [jenni](https://github.com/m-sureshraj/jenni) is a Jenkins Personal Assistant
* [jenkins-job-cli](https://github.com/gocruncher/jenkins-job-cli) Easy way to run Jenkins job from the Command Line

## Stargazers over time

{{printStarHistory "jenkins-zh" "jenkins-cli"}}

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fjenkins-zh%2Fjenkins-cli.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fjenkins-zh%2Fjenkins-cli?ref=badge_large)
