package framework

import (
	"github.com/TeemoKill/BBot/consts"
	"github.com/spf13/viper"
)

func (b *BBot) loadConfig() (err error) {
	b.Cfg = viper.New()

	b.Cfg.SetConfigName(consts.BBotConfigName)
	b.Cfg.SetConfigType(consts.BBotConfigType)
	b.Cfg.AddConfigPath(consts.BBotConfigPath)

	err = b.Cfg.ReadInConfig()
	if err != nil {
		return err
	}

	b.Cfg.WatchConfig()

	return nil
}
