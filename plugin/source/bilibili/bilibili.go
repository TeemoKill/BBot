package bilibili

import (
	"github.com/TeemoKill/BBot/framework"
	"github.com/spf13/viper"
)

type BiliWatcher struct {
	Cfg *viper.Viper

	noticeChan chan *framework.Notification
}

func (b *BiliWatcher) Init(cfg *viper.Viper) (err error) {
	// TODO: implement
	b.Cfg = cfg

	return err
}

func (b *BiliWatcher) Stop() (err error) {
	// TODO: implement

	return err
}

func (b *BiliWatcher) Name() string {
	return BilibiliModuleName
}

// NoticeTo implements framework.Source interface
func (b *BiliWatcher) NoticeTo(c chan *framework.Notification) {
	b.noticeChan = c
}
