---
title: "completion"
weight: 102
---

## install auto-completion for zsh

if you install iterm2 on your macOS or linux，and you use `oh-my-zsh`，you can follow the steps：

```
# cd ~/.oh-my-zsh/plugins
// create incr folder
# mkdir incr
// download incr plugin
# wget https://mimosa-pudica.net/src/incr-0.2.zsh
// authorize incr 
# chmod 777 ~/.oh-my-zsh/plugins/incr/incr-0.2.zsh
# vim ~/.zshrc, and insert "source ~/.oh-my-zsh/plugins/incr/incr-0.2.zsh" in the "~/.zshrc"，save and quit
// flush configuration
# source ~/.zshrc
```

Then you can use auto-completion for jcli, maybe you will find other commands can also use it
