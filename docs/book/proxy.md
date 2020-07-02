---
title: "代理"
weight: 90
---

你可能需要设置代理才可以访问到 Jenkins，这时候，可以给 `jcli` 配置代理服务器的信息。
执行命令：`jcli config edit` 就会打开配置文件，参考下面的配置：

```
jenkins_servers:
- name: dev
  url: http://192.168.1.10
  username: admin
  token: 11132c9ae4b20edbe56ac3e09cb5a3c8c2
  proxy: http://192.168.10.10:47586
  proxyAuth: username:password
```
