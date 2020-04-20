## JFrog Bintray
You can download `jcli` from [bintray.com/jenkins-zh](https://bintray.com/beta/#/jenkins-zh/generic/jenkins-cli/).

`curl -L "https://bintray.com/jenkins-zh/jenkins-cli/download_file?file_path=v0.0.24%2Fjcli-darwin-amd64.tar.gz"|tar xzv`

Get all versions from [here](https://dl.bintray.com/jenkins-zh/generic/jenkins-cli/).

## YUM

Add YUM source repo by the following command:

```shell script
wget https://bintray.com/jenkins-zh/rpm/rpm -O /etc/yum.repos.d/bintray-jcli.repo
```

then you can install it by: `yum install jcli`

## Debian

Add deb source repo by the following command:

```shell script
echo "deb https://dl.bintray.com/jenkins-zh/deb wheezy main" | sudo tee -a /etc/apt/sources.list
```

then you can install it by: `sudo apt-get install jcli`

## Image
Also you can try the following ways:

`jcli_id=$(docker create jenkinszh/jcli) && sudo docker cp $jcli_id:/usr/local/bin/jcli /usr/local/bin/jcli && docker rm -v $jcli_id`

Download different version of OS? Just need to change the docker image tag:

|image|description|
|---|---|
|`jenkinszh/jcli`|Linux|
|`jenkinszh/jcli:darwin`|Mac|
|`jenkinszh/jcli:win`|Windows, you can find it from `/usr/local/bin/jcli.exe`|
|`jenkinszh/jcli:dev`|Developing version, find can find them from `/bin/linux/jcli` or `/bin/darwin/jcli` or `/bin/windows/jcli.exe`|

Want to try the latest features? Download the developing version of different platform:

- `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/linux/jcli . && docker rm -v $jcli_id`
- `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/darwin/jcli . && docker rm -v $jcli_id`
- `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/windows/jcli.exe . && docker rm -v $jcli_id`
