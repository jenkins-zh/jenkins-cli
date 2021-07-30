命令钩子，允许你在执行命令前后，执行特定的命令；钩子包括有前置和后置命令。

例如：我们可以给执行上传插件的命令添加钩子，上传前构建插件项目，上传完成后重启 Jenkins

```
preHooks:
- path: plugin.upload
  cmd: mvn clean package -DskipTests -Dmaven.test.skip
postHooks:
- path: plugin.upload
  cmd: jcli center watch --util-install-complete
- path: plugin.upload
  cmd: jcli restart -b
- path: plugin.upload
  cmd: mvn clean
```

所谓前置钩子也就是 `preHooks`，后置钩子为 `postHooks`。字段 `path` 为以点（.）链接的命令。
其中，钩子命令依照所配置的顺序执行。
