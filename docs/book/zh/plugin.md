---
title: 插件
weight: 70
---

# 插件

`jcli` 可以让你搜索、下载、安装、卸载或者上传插件。

## 列表

下面的命令可以列出所有已经安装的插件：

`jcli plugin list`

如果需要根据字段进行过滤的话，可以利用下面的命令：

`jcli plugin list --filter ShortName=github`

## 检索

你可以通过关键字来搜索要安装的插件，命令如下：

`jcli plugin search zh-cn`

## 安装

给定要安装的插件的名称，并用如下的命令来安装：

`jcli plugin install localization-zh-cn`

执行完成上面的安装命令后，可以通过下面的命令看到安装过程：

`jcli center watch`

## 下载

当你的 Jenkins 无法访问外网，或者其他无法直接安装插件的情况下，可以先把需要安装的插件下载到本地，然后再上传。

`jcli plugin download localization-zh-cn`

默认情况下，会下载你需要的插件以及依赖。如果不需要下载依赖的话，可以使用参数： `--skip-dependency`

## 上传

你可以选择上传本地或者远程的插件文件，甚至可以实现编译本地的插件源码后上传。在没有给定任何参数的情况下，
上传命令首先会尝试执行 Maven 的构建命令，然后再上传文件。

`jcli plugin upload`

如果你已经有编译好的插件文件，可以使用下面的命令：

`jcli plugin upload sample.hpi`

## 升级

如果没有任何参数的话，下面的命令会列出来所有可以升级的插件，利用方向键以及空格可以选择所需要升级的插件，最后回车确认：

`jcli plugin upgrade`

另外，也可以通过给定插件名称的方式，直接升级指定的插件：

`jcli plugin upgrade blueocean-personalization`

## 卸载

`jcli plugin uninstall`

## 检查更新

检查更新，也就是从 Jenkins 的更新中心（Update Center）中获取最新的版本信息，执行下面的命令：

`jcli plugin check`

该命令执行的时间长短，和 Jenkins 所在机器的网络状态有关系，默认的超时时间为：30秒。另外，也可以通过设置参数的方式指定：

`jcli plugin checkout --timeout 60`

## 创建

对于插件的开发者而言，插件的创建、构建、发布也是高频操作，`jcli` 对这些都有支持：

`jcli plugin create`

## 构建

`jcli plugin build`

## 发布

`jcli plugin release`
