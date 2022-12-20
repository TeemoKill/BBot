package framework

import "github.com/spf13/viper"

type Source interface {
	Init(cfg *viper.Viper) error
	Name() string
	NoticeTo(chan *Notification)

	Stop() error
}
