package consts

import "fmt"

const (
	BBotConfigName = "bbot"
	BBotConfigType = "toml"
	BBotConfigPath = "."
)

var (
	BBotConfigFilePath = fmt.Sprintf("%s.%s", BBotConfigName, BBotConfigType)
)
