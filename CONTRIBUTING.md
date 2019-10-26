Welcome! Any kinds of contributions are very welcome. Please go through our contribution
guide before you create a Pull Request for `jcli`.

## CLI

`jcli` is a command line interface. So a CLI framework is super important for us. Thanks to 
[cobra](https://github.com/spf13/cobra). It powers us to do a better job.

## Jenkins REST API

API is another important part of this project. `jcli` manages Jenkins by the REST API.
There is no full specification for this API at the moment, Jenkins core and plugins provide documentation independently.
You can figure it by yourself, or just join our
[gitter room](https://gitter.im/jenkinsci/jenkins-cli) to ask about specific APIs if needed.

Useful links:

* [Jenkins Remote Access API](https://wiki.jenkins.io/display/JENKINS/Remote+access+API)
* [Jenkins REST API overview](https://www.youtube.com/watch?v=D93t1jElt4Q) by [Cliffano Subagio](https://github.com/cliffano)
* 

## Testing

We use a BDD Testing Framework to test our project. Please make sure you're familiar
with [ginkgo](https://github.com/onsi/ginkgo) before you get start to contribute.

## Pull Requests

Before you get started, please fork this project into your GitHub account. Then
create a git branch base on what you want to improve. Please consider **never** using
the master branch as your development branch. And the behaviour of the git **force push** is not
encouraged when submitting pull requests.

Please **do not** create another Pull Request if you messed up your git commit records.

In order to generate nice [release notes](https://github.com/jenkins-zh/jenkins-cli/releases),
please consider writing a proper Pull Request title.
[release-draft](https://github.com/toolmantim/release-drafter) will generate the notes base your title.

## Quality

Quality is the heart of a project. So please make sure your Pull Request could pass the
[Sonar Quality Gate](https://sonarcloud.io/dashboard?id=jenkins-zh_jenkins-cli).

|Metric|Operator|Value|
|---|---|---|
|Coverage|is less than|90.0%|
|Duplicated Lines(%)|is greater than|3.0%|
|Maintainablity Rating|is worse than|A|
|Blocker Issues|is greater than|1|
|Code Smells|is greater than|1|
|Reliablity Rating|is worse than|A|
|Security Rating|is worse than|A|

## Good Start

The [newbie](https://github.com/jenkins-zh/jenkins-cli/issues?q=is%3Aissue+is%3Aopen+label%3Anewbie) issues
are the good start.

## Git Backup

We use [git-backup-actions](https://github.com/jenkins-zh/git-backup-actions/) to backup this repo into 
[gitee](https://gitee.com/jenkins-zh/jenkins-cli).
