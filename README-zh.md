# Jenkins CLI

[![Go Report Card](https://goreportcard.com/badge/jenkins-zh/jenkins-cli)](https://goreportcard.com/report/jenkins-zh/jenkins-cli)

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

该工具还处在非常早期的开发阶段，我们欢迎任何形式的贡献。
