创建一个子 Shell 并执行 `jcli` 命令。这时候，不管如何修改 `jcli` 的配置文件，退出后都不会影响之前的配置。

执行下面的命令，会将所选择的 Jenkins 配置 `local` 作为默认的值：

`jcli shell local`

此时，我们执行命令 `jcli config` 的话能看出来。
