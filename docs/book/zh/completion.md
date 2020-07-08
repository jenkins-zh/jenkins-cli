---
title: 自动补全
weight: 102
---

# 自动补全

如果你已经在 mac 或 linux 上使用的是 `oh-my-zsh`，你可以尝试以下步骤：

```text
# cd ~/.oh-my-zsh/plugins
// 创建 incr 文件夹
# mkdir incr
// 下载 incr 插件
# wget https://mimosa-pudica.net/src/incr-0.2.zsh
// 对 incr 进行授权
# chmod 777 ~/.oh-my-zsh/plugins/incr/incr-0.2.zsh
# vim ~/.zshrc,然后在 “~/.zshrc” 文件中加入 “source ~/.oh-my-zsh/plugins/incr/incr-0.2.zsh”，保存退出
// 更新配置
# source ~/.zshrc
```

接下来，就可以使用 jcli 的自动补全功能了，而且你可能发现不仅仅只有 jcli 可以自动补全，很多命令都可以自动补全了

