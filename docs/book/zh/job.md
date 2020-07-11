---
title: 任务
weight: 80
---

## 搜索

使用如下的命令可以搜索 Jenkins 任务：

`jcli job search input`

要查找特定类型的 Jenkins 任务，可以通过过滤对应字段的值来实现。下面，给出一个查找参数化任务的例子：

`jcli job search --filter Parameterized=true`

其中 `--filter` 支持任意字段，它是以是否包含指定字符串进行判断的。

## 构建

要触发一个任务的话，可以使用下面的命令：

`jcli job build "jobName" -b`

## 交互式输入

执行到 Jenkins 流水线中的 `input` 指令时，会有交互式输入的提示。下面是一个样例：

```
pipeline {
    agent {
        label 'master'
    }
    
    stages {
        stage('sample') {
            steps {
                input 'test'
            }
        }
    }
}
```

运行上面的流水线后，执行到 `input` 位置就会阻塞并等待输入，此时可以通过命令 `jcli job input test 1` 来使得继续执行或者中断。

## 编辑

目前，只对以脚本的形式保存在 Jenkins 上的流水线有编辑功能的支持。命令非常简单：`jcli job edit test`

如果希望能快速地给出一个流水线的样例的话，当在流水线脚本为空时，可以执行命令：`jcli job edit test --sample`

如果希望编辑流水线并保存退出后，直接触发的话，可以使用对应的参数来实现：`jcli job edit test --build`

## 禁用

禁用任务：`jcli job disable job/test/`

启用任务：`jcli job enable job/test/`

## 查看日志

通过下面的命令可以参考一个任务的执行日志：

`jcli job log "jobName" -w`

## 查看历史

`jcli job history job/test/`

## 归档文件

查看归档文件列表 `jcli job artifact job/test/`

下载归档文件 `job artifact download /job/tsf/job/ddd/`

## 显示指定列

当以表格形式输出，希望能输出指定的字段为列时，我们可以通过下面的方式实现：

`jcli job search --columns Name,URL,Parameterized`

请注意，上面的参数 `--columns` 的值是以英文逗号（,）为分割的。

如果不希望输出表头，可以增加参数：`--no-headers`

`jcli job search --columns Name,URL,Parameterized --no-headers`

## 任务类型

列出当前 Jenkins 所支持的任务类型 `jcli job type`
