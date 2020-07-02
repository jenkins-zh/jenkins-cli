---
title: 任务
weight: 80
---

# 任务

使用如下的命令可以搜索 Jenkins 任务：

`jcli job search input`

要触发一个任务的话，可以使用下面的命令：

`jcli job build "jobName" -b`

通过下面的命令可以参考一个任务的执行日志：

`jcli job log "jobName" -w`

## 搜索任务

要查找特定类型的 Jenkins 任务，可以通过过滤对应字段的值来实现。下面，给出一个查找参数化任务的例子：

`jcli job search --filter Parameterized=true`

其中 `--filter` 支持任意字段，它是以是否包含指定字符串进行判断的。

## 显示指定列

当以表格形式输出，希望能输出指定的字段为列时，我们可以通过下面的方式实现：

`jcli job search --columns Name,URL,Parameterized`

请注意，上面的参数 `--columns` 的值是以英文逗号（,）为分割的。

如果不希望输出表头，可以增加参数：`--no-headers`

`jcli job search --columns Name,URL,Parameterized --no-headers`

