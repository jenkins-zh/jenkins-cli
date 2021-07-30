---
title: 凭据
weight: 101
since: v0.0.24
---

# 凭据

通过 `jcli` 可以在 Jenkins 上创建凭据（Credentials），下面介绍使用方法。

## 创建

Jenkins 中的凭据有多种类型，下面的命令会创建一个用户名和密码类型的凭据：

```text
jcli credential create --credential-username your-username \
--credential-password your-password --desc your-credential-remark
```

下面的命令创建一个只包含单一加密文本的凭据：

`jcli credential create --secret my-secret --type secret`

## 列表

`jcli credential list`

## 删除

我们可以根据 Jenkins 凭据的唯一标示来删除：

`jcli credential delete --id b0b0f865-f0c0-477c-a5ba-9fae88477f9e`

