# Jenkins CLI

[![Go Report Card](https://goreportcard.com/badge/jenkins-zh/jenkins-cli)](https://goreportcard.com/report/jenkins-zh/jenkins-cli)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=jenkins-zh_jenkins-cli&metric=alert_status)](https://sonarcloud.io/dashboard?id=jenkins-zh_jenkins-cli)
![Sonar Coverage](https://img.shields.io/sonar/coverage/jenkins-zh_jenkins-cli?server=https%3A%2F%2Fsonarcloud.io)
[![Travis](https://img.shields.io/travis/jenkins-zh/jenkins-cli.svg?logo=travis&label=build&logoColor=white)](https://travis-ci.org/jenkins-zh/jenkins-cli)
[![Contributors](https://img.shields.io/github/contributors/jenkins-zh/jenkins-cli.svg)](https://github.com/jenkins-zh/jenkins-cli/graphs/contributors)
[![GitHub release](https://img.shields.io/github/release/jenkins-zh/jenkins-cli.svg?label=release)](https://github.com/jenkins-zh/jenkins-cli/releases/latest)
![GitHub All Releases](https://img.shields.io/github/downloads/jenkins-zh/jenkins-cli/total)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/jenkins-zh/jenkins-cli)
[![Gitter](https://badges.gitter.im/jenkinsci/jenkins-cli.svg)](https://gitter.im/jenkinsci/jenkins-cli?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

Jenkins CLI 可以帮忙你轻松地管理 Jenkins。不管你是一名插件开发者、管理员或者只是一个普通的 Jenkins 用户，它都是为你而生的！

# 特点

* 支持多 Jenkins 实例管理
* 插件管理（查看列表、搜索、安装、上传）
* 任务管理（搜索、构建触发、日志查看）
* 在浏览器中打开你的 Jenkins
* 重启你的 Jenkins
* 支持通过代理连接

# 安装

我们目前支持的操作系统包括：MacOS、Linux 以及 Widnows。

## mac

你可以通过 `brew` 来安装 jcli。
```
brew tap jenkins-zh/jcli
brew install jcli
```

## Linux

要在 Linux 操作系统上安装 `jcli` 的话，非常简单。只需要执行下面的命令即可：
```
curl -L https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-linux-amd64.tar.gz|tar xzv
sudo mv jcli /usr/local/bin/
```

## Windows

你只要[点击这里](https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-windows-386.tar.gz)就可以下载到最新版本的压缩包。之后，把解压后的文件 `jcli` 拷贝到你的系统目录下即可。

## 其他包管理器

这里还有一些其他的包管理器：

* [GoFish](https://gofi.sh/) 的用户可以使用命令 `gofish install jcli` 来安装。

# 入门

当安装 `jcli` 以后。你需要提供一份配置文件。请执行命令 `jcli config generate`，该命令会帮助你编辑配置文件 `~/.jenkins-cli.yaml` ，你需要根据实际的 Jenkins 配置情况做相应的修改。

# 贡献

如果你对该项目感兴趣，请首先仔细阅读我们的[贡献指南](CONTRIBUTING.md)。我们欢迎任何形式的贡献。
