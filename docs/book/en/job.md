---
title: "Job"
weight: 80
---

## Job Management

You can search a job list using a keyword, like this: `jcli job search input`.

It's very simple to trigger a job. We have the batch mode and interactive mode. This command will finish immediately.

`jcli job build "folderName jobName" -b`

Once you triggered a job, then you can watch the log output by `jcli job log "zjproject zjproject-inputstep55" -w`. This command will always output the log of the last build.

## Proxy Support

Jenkins might be stay in behind a firewall. So we cannot connect it directly. You can give `jcli` a proxy setting. It's also very simple to support a proxy setting. You just need to execute: `jcli config edit`. Then find the item which you want to add a proxy. Like the below demo:

```
jenkins_servers:
- name: dev
  url: http://192.168.1.10
  username: admin
  token: 11132c9ae4b20edbe56ac3e09cb5a3c8c2
  proxy: http://192.168.10.10:47586
  proxyAuth: username:password
```
