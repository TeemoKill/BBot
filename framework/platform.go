package framework

import "github.com/spf13/viper"

type Platform interface {
	Init(cfg *viper.Viper) error
	Name() string

	Stop() error
}
