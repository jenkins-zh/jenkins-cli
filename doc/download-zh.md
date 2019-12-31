## JFrog Bintray
你可以从 [bintray.com/jenkins-zh](https://bintray.com/jenkins-zh/jenkins-cli/jenkins-cli) 下载 `jcli`.

`curl -L "https://bintray.com/jenkins-zh/jenkins-cli/download_file?file_path=jcli-darwin-amd64.tar.gz"|tar xzv`

## 镜像
如果您无法从 GitHub 上下载 `jcli`，请尝试下面的方法：

`jcli_id=$(docker create jenkinszh/jcli) && sudo docker cp $jcli_id:/usr/local/bin/jcli /usr/local/bin/jcli && docker rm -v $jcli_id`

要下载不同操作系统下的二进制文件？只需要修改 docker 容器的标签：

|镜像|描述|
|---|---|
|`jenkinszh/jcli`|Linux|
|`jenkinszh/jcli:darwin`|Mac|
|`jenkinszh/jcli:win`|Windows，你可以从 `/usr/local/bin/jcli.exe` 这里找到可执行程序|
|`jenkinszh/jcli:dev`|你可以从这里找到开发版本 `/bin/linux/jcli` 、`/bin/darwin/jcli` 或 `/bin/windows/jcli.exe`|

想要体验最新的特性？您可以下载不同平台下的开发版本：

- `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/linux/jcli . && docker rm -v $jcli_id`
- `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/darwin/jcli . && docker rm -v $jcli_id`
- `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/windows/jcli.exe . && docker rm -v $jcli_id`
