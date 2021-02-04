# download

## YUM

Add YUM source repo by the following command:

```
cat > bintray-jenkins-zh-rpm.repo <<EOF
#bintraybintray-jenkins-zh-rpm - packages by jenkins-zh from Bintray
[bintraybintray-jenkins-zh-rpm]
name=bintray-jenkins-zh-rpm
baseurl=https://dl.bintray.com/jenkins-zh/rpm
gpgcheck=0
repo_gpgcheck=0
enabled=1
EOF
sudo mv bintray-jenkins-zh-rpm.repo /etc/yum.repos.d/
sudo yum update
```

then you can install it by: `yum install jcli`

## Debian

Add deb source repo by the following command:

```
echo "deb [trusted=yes] https://dl.bintray.com/jenkins-zh/deb wheezy main" | sudo tee -a /etc/apt/sources.list
sudo apt update
```

then you can install it by: `sudo apt-get install jcli`

## Image

Also you can try the following ways:

`jcli_id=$(docker create jenkinszh/jcli) && sudo docker cp $jcli_id:/usr/local/bin/jcli /usr/local/bin/jcli && docker rm -v $jcli_id`

Download different version of OS? Just need to change the docker image tag:

| image | description |
| :--- | :--- |
| `jenkinszh/jcli` | Linux |
| `jenkinszh/jcli:darwin` | Mac |
| `jenkinszh/jcli:win` | Windows, you can find it from `/usr/local/bin/jcli.exe` |
| `jenkinszh/jcli:dev` | Developing version, find can find them from `/bin/linux/jcli` or `/bin/darwin/jcli` or `/bin/windows/jcli.exe` |

Want to try the latest features? Download the developing version of different platform:

* `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/linux/jcli . && docker rm -v $jcli_id`
* `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/darwin/jcli . && docker rm -v $jcli_id`
* `jcli_id=$(docker create jenkinszh/jcli:dev) && sudo docker cp $jcli_id:/bin/windows/jcli.exe . && docker rm -v $jcli_id`

## Out-of-date

Below distributions are out-of-date. If you want to maintain them, please let us know.

* [GoFish](https://gofi.sh/) users can use `gofish install jcli`
* [Chocolatey](https://chocolatey.org/packages/jcli) users can use `choco install jcli`
* [Snapcraft](https://snapcraft.io/jcli) users can use `sudo snap install jcli`

