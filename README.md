# BBot


`BBot` 是一个消息推送框架，内置支持B 站[直播/动态] 消息的关注  
内置了基于 [MiraiGo](https://github.com/Mrs4s/MiraiGo) 的QQ 消息推送插件  
也支持通过自定义插件扩展订阅源或推送目标。


## 设计目标

- 易于扩展
- 结构清晰
- 方便配置订阅源和推送
- 方便的命令权限、频率控制
- 高性能
- 兼容`DDBOT` 数据库(buntdb)
    - 存储使用接口实现(kv 接口), 方便扩展其他种类的存储
- 低资源消耗(暂时优先级不高)


## 参考项目

go-cqhttp  
https://github.com/Mrs4s/go-cqhttp

DDBOT  
https://github.com/Sora233/DDBOT

MiraiValBot  
https://github.com/eric2788/MiraiValBot


## 声明

- 您可以免费使用`BBot` 进行其他商业活动，但不允许通过出租、出售`BBot` 等方式进行商业活动。
- 如果您部署了私人的bot 实例，可以接受他人对您私人部署的bot 进行捐赠以帮助bot 运行，但捐赠必须本着自愿的原则，不允许用bot 使用权来强制或变相强制他人进行捐赠。
- 如果您使用了`BBot` 的源代码，或者对`BBot` 源代码进行修改，您应该用相同的开源许可（AGPL3.0）进行开源，并标明著作权。


## BBot :star:趋势图

[![Stargazers over time](https://starchart.cc/TeemoKill/BBot.svg)](https://starchart.cc/TeemoKill/BBot)
