---
title: 插件
weight: 70
---

# 插件

`jcli` 可以让你搜索、下载、安装、卸载或者上传插件。

## 安装插件

你可以通过关键字来搜索要安装的插件，命令如下：

`jcli plugin search zh-cn`

然后，拷贝要安装的插件的名称，并用如下的命令来安装：

`jcli plugin install localization-zh-cn`

## 下载插件

当你的 Jenkins 无法访问外网，或者其他无法直接安装插件的情况下， 可以先把需要安装的插件下载到本地，然后再上传。

`jcli plugin download localization-zh-cn`

默认情况下，会下载你需要的插件以及依赖。如果不需要下载依赖的话，可以使用参数： `--skip-dependency`

## 上传插件

你可以选择上传本地或者远程的插件，甚至可以实现编译本地的插件源码后上传。 在没有给定任何参数的情况下，上传命令首先会尝试执行 Maven 的构建命令， 然后再上传文件。

`jcli plugin upload`

如果你已经有编译好的插件文件，可以使用下面的命令：

`jcli plugin upload sample.hpi`

