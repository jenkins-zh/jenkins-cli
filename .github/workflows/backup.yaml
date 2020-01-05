name: Backup Git repository

on:
  push:
    branches:
    - master

jobs:
  hugo-deploy-gh-pages:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: backup
      uses: jenkins-zh/git-backup-actions@v0.0.1
      env:
        GIT_DEPLOY_KEY: ${{ secrets.GIT_DEPLOY_KEY }}
        TARGET_GIT: "git@gitee.com:jenkins-zh/jenkins-cli.git"
