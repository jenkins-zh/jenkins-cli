# 计算节点

Jenkins 的最佳实践是让 master 只做调度任务，其他的构建等任务的执行都放在 agent（计算节点）上运行。
在安装不同插件后，使得 Jenkins 可以支持静态、动态类型的节点。所谓静态，指的是需要我们人工来维护，例如：
创建、上线、下线对应的节点。所谓动态，则可以根据既定的规则，自动地创建、销毁节点；
以 [Kubernetes 插件](https://github.com/jenkinsci/kubernetes-plugin/) 为例，它通过动态地创建
和销毁 [Pod](https://kubernetes.io/docs/concepts/workloads/pods/pod/) 来提供节点的运行。

## 协议

不管是动态还是静态的节点，都需要特定的协议来链接 agent 和 master。Jenkins 可以通过以下协议建立链接：
* SSH
* [JNLP](https://docs.oracle.com/javase/tutorial/deployment/deploymentInDepth/jnlp.html)
* [WMI](https://en.wikipedia.org/wiki/Windows_Management_Instrumentation)

查看节点列表：`jcli agent list`

## 静态节点

```
jcli agent create macos
jcli agent launch macos
```

当前，只支持 JNLP 类型的节点创建。另外，对于需要通过 HTTP 代理才能链接到 Jenkins 的话，暂时不支持。

## 删除节点

给定节点的名称即可删除：`jcli agent delete macos`
