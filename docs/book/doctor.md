---
title: 诊断
weight: 900
since: v0.0.24
---

# 诊断

由于错误配置或者是缺少相应插件，可能会导致 `jcli` 无法正常工作。然而，有时候想要快速地找到问题所在， 是一件不容易而且费时的事情。这里要介绍的`诊断`功能，就是为了解决这样的问题而存在的。

## 插件依赖

就像命令 `jcli job search` 要依赖插件 `[pipeline-restful-api](https://plugins.jenkins.io/pipeline-restful-api)` 一样，其他部分插件也有类似的依赖。有的情况下，还对插件的版本有要求。

在执行命令时，如果发现无法使用，可以尝试使用诊断参数来检查是否缺少依赖：

`jcli job search --doctor`

其中 `--doctor` 是一个全局参数。当有依赖不满足等情况发生时，会有相应的错误提示信息输出。例如： `Error: lack of plugin pipeline-restful-api`。

