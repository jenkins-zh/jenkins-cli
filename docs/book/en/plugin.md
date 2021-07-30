---
title: "Plugin"
anchor: "plugin"
weight: 70
---

## Plugin Management

`jcli` allows you to search, download, install, uninstall or upload a plugin.

First, please search a plugin by a keyword if you want to install it. You can get a plugin list by execute `jcli plugin search zh-cn`. You can install it with the plugin name.

For example, you can install the Simplified Chinese Localization plugin by `jcli plugin install localization-zh-cn`.

### Download Plugins

Some times, Jenkins just cannot connect with the offical Update Center. We can use the `download` sub-cmd to download all the plugins which're you need, then upload them. This command will take care of the dependencies of the plugin.

You can try it:

`jcli plugin download localization-zh-cn`
