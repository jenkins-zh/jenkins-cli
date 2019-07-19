# Jenkins CLI

Jenkins CLI allows you manage your Jenkins as an easy way. No matter you're a plugin
developer, administrator or just a regular user, it borns for you!

# Features

* Multiple Jenkins support
* Plugins management (list, search, install, upload)
* Job management (search, build, log)
* Open your Jenkins with a browse
* Restart your Jenkins
* Connection with proxy support

# Get started

We support mac, linux and windows for now.

## mac

You can use `brew` to install jcli.
```
brew tap linuxsuren/jcli
brew install jcli
```

## Linux

It's very simple to install `jcli` into your Linux OS. Just need to execute a command line at below:
```
ostype=linux-amd64&&curl -s https://api.github.com/repos/LinuxSuRen/jenkins-cli/releases/latest | grep -e 'browser_download_url.*jcli-'${ostype}'.tar.gz'|awk '{print $2}'|xargs wget&&tar xzvf jcli-${ostype}.tar.gz
sudo mv jcli /usr/local/bin/&&rm jcli-${ostype}.tar.gz
```

## Windows

You can find the right version from the [release page](https://github.com/LinuxSuRen/jenkins-cli/releases). Then download the tar file, cp the uncompressed `jcli` directory into your system path.

# Contribution

It's still under very early develope time. Any contribution is welcome.
