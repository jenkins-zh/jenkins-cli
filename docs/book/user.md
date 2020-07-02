---
title: "用户"
weight: 100
---

`jcli` 可以完成用户的创建、删除以及生成令牌（Token）的操作，

## 创建用户

```
jcli user create <username> [password] [flags]
```

在创建用户的时候，可以指定一个密码或者随机生成。

## 生成令牌

Jenkins 的 Web API 必须是通过令牌（Token）来访问，`jcli` 支持给当前用户或者
指定用户生成令牌。给当前用户生成令牌的命令如下：

`jcli user token -g`

如果希望通过管理员给其他的 Jenkins 用户生成令牌的话，需要在启动 Jenkins 时给定一些参数，
具体参考下面的命令：

```
jcli center start --admin-can-generate-new-tokens
jcli user token -g --target-user target-user-name
```

上面的第一条命令会启动 Jenkins 并设置为允许有管理员权限的用户为其他用户生成令牌。
