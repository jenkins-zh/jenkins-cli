# 快速开始

[English](https://github.com/jenkins-zh/jenkins-cli/blob/master/README.md)

## Jenkins CLI

[![](https://goreportcard.com/badge/jenkins-zh/jenkins-cli)](https://goreportcard.com/report/jenkins-zh/jenkins-cli) [![](https://sonarcloud.io/api/project_badges/measure?project=jenkins-zh_jenkins-cli&metric=alert_status)](https://sonarcloud.io/dashboard?id=jenkins-zh_jenkins-cli) [![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/jenkins-zh/jenkins-cli) ![Sonar Coverage](https://img.shields.io/sonar/coverage/jenkins-zh_jenkins-cli?server=https%3A%2F%2Fsonarcloud.io) [![Travis](https://img.shields.io/travis/jenkins-zh/jenkins-cli.svg?logo=travis&label=build&logoColor=white)](https://travis-ci.org/jenkins-zh/jenkins-cli) [![Contributors](https://img.shields.io/github/contributors/jenkins-zh/jenkins-cli.svg)](https://github.com/jenkins-zh/jenkins-cli/graphs/contributors) [![GitHub release](https://img.shields.io/github/release/jenkins-zh/jenkins-cli.svg?label=release)](https://github.com/jenkins-zh/jenkins-cli/releases/latest) ![GitHub All Releases](https://img.shields.io/github/downloads/jenkins-zh/jenkins-cli/total) [![Docker Pulls](https://img.shields.io/docker/pulls/jenkinszh/jcli.svg)](https://hub.docker.com/r/jenkinszh/jcli/tags) ![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/jenkins-zh/jenkins-cli) [![Gitter](https://badges.gitter.im/jenkinsci/jenkins-cli.svg)](https://gitter.im/jenkinsci/jenkins-cli?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) [![HitCount](http://hits.dwyl.com/jenkins-zh/jenkins-cli.svg)](http://hits.dwyl.com/jenkins-zh/jenkins-cli)

Jenkins CLI 可以帮忙你轻松地管理 Jenkins。不管你是一名插件开发者、管理员或者只是一个普通的 Jenkins 用户，它都是为你而生的！

## 特性

* 支持多 Jenkins 实例管理
* 插件管理（查看列表、搜索、安装、上传）
* 任务管理（搜索、构建触发、日志查看）
* 支持配置即管理
* 在浏览器中打开你的 Jenkins
* 重启你的 Jenkins
* 支持通过代理连接

## 安装

我们目前支持的操作系统包括：MacOS、Linux 以及 Windows。

### mac

你可以通过 `brew` 来安装 jcli。

```text
brew tap jenkins-zh/jcli
brew install jcli
```

### Linux

要在 Linux 操作系统上安装 `jcli` 的话，非常简单。只需要执行下面的命令即可：

```text
curl -L https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-linux-amd64.tar.gz|tar xzv
sudo mv jcli /usr/local/bin/
```

### Windows

你只要[点击这里](https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-windows-amd64.zip)就可以下载到最新版本的压缩包。之后，把解压后的文件 `jcli` 拷贝到你的系统目录下即可。

### 其他包管理器

这里还有一些其他的包管理器：

* [Scoop](https://scoop.sh/) 的用户可以使用命令 `scoop install jcli` 来安装

了解更多[如何下载 jcli](https://github.com/jenkins-zh/jenkins-cli/tree/e83af606f648040665b8b2955c1c2414bb68c1db/docs/book/zh/download-zh.md). 你可以从[这里](http://somsubhra.com/github-release-stats/?username=jenkins-zh&repository=jenkins-cli)获取下载的统计信息。

## 入门

查阅[官方文档](http://jcli.jenkins-zh.cn/)可以了解到更多有关如何使用 `jcli` 的信息。

或者，你可以选择 Jenkins CLI 的[一个在线的交互式教程](https://www.katacoda.com/jenkins-zh/scenarios/course-jcli)。

## 插件

通过安装插件可以增强 Jenkins CLI 的功能。按照下面的命令就可以安装一个插件：

```text
jcli config plugin fetch
jcli config plugin install account
```

所有官方的插件，都可以在[这里](https://github.com/jenkins-zh/jcli-plugins)找到。

## 贡献

如果你对该项目感兴趣，请首先仔细阅读我们的[贡献指南](https://github.com/jenkins-zh/jenkins-cli/tree/e83af606f648040665b8b2955c1c2414bb68c1db/CONTRIBUTING.md)。我们欢迎任何形式的贡献。

感谢 JetBrains 为我们提供了开源许可证。  
[![goland.svg](../.gitbook/assets/goland%20%282%29.svg)](https://www.jetbrains.com/?from=jenkins-cli)

## 相关的项目

有一些相关的项目你可能会比较感兴趣：

* [jenni](https://github.com/m-sureshraj/jenni) 是一个 Jenkins 个人助手

## 点赞数趋势图

[![Stargazers over time](https://starchart.cc/jenkins-zh/jenkins-cli.svg)](https://starchart.cc/jenkins-zh/jenkins-cli)

