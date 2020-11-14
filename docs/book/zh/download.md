# 下载

## JFrog Bintray

你可以从 [bintray.com/jenkins-zh](https://bintray.com/beta/#/jenkins-zh/generic/jenkins-cli/) 下载 `jcli`.

`curl -L "https://bintray.com/jenkins-zh/jenkins-cli/download_file?file_path=v0.0.24%2Fjcli-darwin-amd64.tar.gz"|tar xzv`

点击[这里](https://dl.bintray.com/jenkins-zh/generic/jenkins-cli/)查看所有版本。

## YUM

通过下面的命令添加 YUM 源：

\`\`\`shell script wget [https://bintray.com/jenkins-zh/rpm/rpm](https://bintray.com/jenkins-zh/rpm/rpm) -O /etc/yum.repos.d/bintray-jcli.repo

```text
然后，你就可以安装了：`yum install jcli`

## Debian

通过下面的命令添加 deb 源：

```shell script
echo "deb https://dl.bintray.com/jenkins-zh/deb wheezy main" | sudo tee -a /etc/apt/sources.list
```

然后，你就可以安装了：`sudo apt-get install jcli`

## 镜像

你也可以尝试下面的方法：

`jcli_id=$(docker create jenkinszh/jcli) && sudo docker cp $jcli_id:/usr/local/bin/jcli /usr/local/bin/jcli && docker rm -v $jcli_id`

要下载不同操作系统下的二进制文件？只需要修改 docker 容器的标签：

| 镜像 | 描述 |
| :--- | :--- |
| `jenkinszh/jcli` | Linux |
| `jenkinszh/jcli:darwin` | Mac |
| `jenkinszh/jcli:win` | Windows，你可以从 `/usr/local/bin/jcli.exe` 这里找到可执行程序 |
| `jenkinszh/jcli:dev` | 你可以从这里找到开发版本 `/bin/linux/jcli` 、`/bin/darwin/jcli` 或 `/bin/windows/jcli.exe` |

想要体验最新的特性？您可以下载不同平台下的开发版本：

* `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/linux/jcli . && docker rm -v $jcli_id`
* `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/darwin/jcli . && docker rm -v $jcli_id`
* `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/windows/jcli.exe . && docker rm -v $jcli_id`

## 过时的

下面的发型版不会及时更新，如果您有兴趣帮忙维护它们的话，请告诉我们，谢谢。

* [GoFish](https://gofi.sh/) 的用户可以使用命令 `gofish install jcli` 来安装
* [Chocolatey](https://chocolatey.org/packages/jcli) 的用户可以使用命令 `choco install jcli` 来安装
* [Snapcraft](https://snapcraft.io/jcli) 的用户可以使用命令 `sudo snap install jcli` 来安装

