package utils

import (
	"runtime"
	"strings"
)

func ExampleConfig() string {
	c := `
log_level = "info"

[qq]
account  = 0 # 你bot 的qq 号，不填则使用扫码登陆
password = "" # 你bot 的qq 密码
  [qq.on_join_group]
  rename = "【BBot】" # bot 进群后自动改名，默认改名为“【BBot】”，如果留空则不自动改名

# 初次运行时将不使用b 站帐号方便进行测试
# 如果不使用b 站帐号，则推荐订阅数不要超过5 个，否则推送延迟将上升
# b 站相关的功能推荐配置一个b 站账号，建议使用小号
# bot 将使用您b 站帐号的以下功能：
#   关注用户 / 取消关注用户 / 查看关注列表
# 请注意，订阅一个账号后，此处使用的b 站账号将自动关注该账号
[bilibili]
SESSDATA = "" # 你的b 站cookie
bili_jct = "" # 你的b 站cookie
interval = 25

[concern]
emit_interval = 5
`
	// 解决windows 上用记事本打开不会正确换行
	if runtime.GOOS == "windows" {
		c = strings.ReplaceAll(c, "\n", "\r\n")
	}
	return c
}
