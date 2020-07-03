# 打开浏览器

我们可以通过下面的命令快速地用浏览器打开 Jenkins：

`jcli open`

## 浏览器设置

默认，`jcli` 会使用系统的缺省浏览器打开。但是，如果希望能用指定的浏览器打开的话，可以参考下面的命令：

```
jcli open --browser "Google-Chrome"
JCLI_BROWSER="Google Chrome" jcli open
```

也就是说，可以通过给定参数，或者设置环境变量的方式来指定浏览器。

## 其他地址

为了方便在浏览器中打开和某个 Jenkins 相关的服务，可以把服务地址添加到配置文件中，例如：

```
current: local
jenkins_servers:
- name: local
  url: http://localhost:8080
  username: admin
  token: '******'
  data:
    baidu: https://baidu.com
    jenkins: https://jenkins.io
```

从上面的配置例子中能看到，字段 `data` 下添加了两个 `key-value`。如果要打开其中的一个地址的话，可以执行下面的命令：

`jcli open .baidu`
