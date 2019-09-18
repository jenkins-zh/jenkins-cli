Welcome! Any kinds of contributions are very welcomed. Please go through our contribution
guide before you try to create a Pull Request for `jcli`.

## CLI

`jcli` is a command line interface. So a CLI framework is super important for us. Thanks to 
[cobra](https://github.com/spf13/cobra). It powers us to do a better job.

## Jenkins API

API is another important part of this project. `jcli` manages your Jenkins by the HTTP API.
There's no official documents for this. You can figure it by yourself, or just join our
[gitter room](https://gitter.im/jenkinsci/jenkins-cli).

## Testing

We use a BDD Testing Framework to test our project. Please make sure you're familar
with [ginkgo](https://github.com/onsi/ginkgo) before you get start to contribute.

## Pull Requests

Before you get start, please fork this project into your GitHub account firstly. Then
create a git branch base on what you want to improve. Please consider **never** using
the master branch as your develope branch. And the behaviour of the git **force push** is not
encourage.

Please **don't** create another Pull Request if you messed up your git commit records.

In order to generate a nice [release notes](https://github.com/jenkins-zh/jenkins-cli/releases),
please consider writing a proper Pull Request title.
[release-draft](https://github.com/toolmantim/release-drafter) will generate the notes base your title.

## Qulity

Qulity is the heart of a project. So please make sure your Pull Request could pass the
[Sonar Qulity Gate](https://sonarcloud.io/dashboard?id=jenkins-zh_jenkins-cli).

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
