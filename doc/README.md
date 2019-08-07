# jcli Document

`jcli` was written by [golang](https://github.com/golang) which could provide you a easy way to manage your Jenkins. Unlike the [build-in cli](https://jenkins.io/doc/book/managing/cli/), `jcli` allows you manage multiple servers.

## How to get it

Read [here](../README.md) to get know about how to install `jcli`.

## Configuration

Once you'v installed `jcli`. You should provide a config file for it. Please execute cmd `jcli config generate`, this will help you to edit the config file `~/.jenkins-cli.yaml`. According to your Jenkins configuration to modify this file.

If you want to modify your config file of `jcli`. You just need to execute `jcli config edit`.

It's simple to add another Jenkins config item. Here's a sample cmd: `jcli config add -n yourJenkinsName --url http://localhost:8080/jenkins --token replacethesampletoken`.

## Plugin Management

`jcli` allows you to search, download, install, uninstall or upload a plugin.

First, please search a plugin by a keyword if you want to install it. You can get a plugin list by execute `jcli plugin search zh-cn`. You can install it with the plugin name.

For example, you can install the Simplified Chinese Localization plugin by `jcli plugin install localization-zh-cn`.

### Download Plugins

Some times, Jenkins just cannot connect with the offical Update Center. We can use the `download` sub-cmd to download all the plugins which're you need, then upload them. This command will take care of the dependencies of the plugin.

You can try it:

`jcli plugin download localization-zh-cn`

## Job Management

TODO.

## Proxy Support

TODO.