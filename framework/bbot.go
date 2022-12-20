package framework

import (
	"github.com/TeemoKill/BBot/consts"
	"github.com/spf13/viper"
	"os"
	"sync"

	"github.com/TeemoKill/BBot/log"
	"github.com/TeemoKill/BBot/utils"
)

type BBot struct {
	Cfg *viper.Viper

	Source   []Source
	Platform []Platform

	noticeChannel chan *Notification
}

func New() *BBot {
	bBot := &BBot{}
	return bBot
}

func (b *BBot) Init() (err error) {
	logger := log.CurrentModuleLogger()

	fileInfo, err := os.Stat(consts.BBotConfigFilePath)
	switch {
	case os.IsNotExist(err):
		logger.WithField("config_filename", consts.BBotConfigFilePath).
			Warnf("没有检测到配置文件，正在生成。如果是第一次运行，可忽略此警告")
		err = b.GenerateExampleConfig()
		if err != nil {
			logger.WithError(err).
				Errorf("failed to generate example config, quitting")
			return err
		}
		return consts.ExampleConfigGeneratedErr

	case err != nil:
		logger.WithError(err).
			WithField("config_filename", consts.BBotConfigFilePath).
			Errorf("检查配置文件失败，退出")
		return err

	default:
		if fileInfo.IsDir() {
			logger.WithField("config_path", consts.BBotConfigFilePath).
				Errorf("检测到配置文件路径，但该路径是目录而不是文件。通常删除该目录并重启bot 可生成新的示例配置，操作前请手动确认该目录下是否有重要文件")
			return consts.ConfigPathNotFileErr
		} else {
			logger.WithField("config_filename", consts.BBotConfigFilePath).
				Infof("检测到配置文件，使用存在的配置文件启动bot")
		}
	}

	err = b.loadConfig()
	if err != nil {
		logger.WithError(err).
			Errorf("读取配置文件失败！请检查配置文件格式是否正确")
		return err
	}

	// register modules

	for _, s := range b.Source {
		err = s.Init(b.Cfg)
		if err != nil {
			logger.WithError(err).
				WithField("source", s.Name()).
				Warnf("source init failed")
		}
	}

	for _, p := range b.Platform {
		err = p.Init(b.Cfg)
		if err != nil {
			logger.WithError(err).
				WithField("platform", p.Name()).
				Warnf("platform init failed")
		}
	}

	return nil
}

func (b *BBot) StartService() {
	logger := log.CurrentModuleLogger()

	// TODO: implement

	logger.Infof("这个BBot 启动完成了")
	logger.Infof("D宝，一款真正人性化的单推BOT")
	/*
		if len(l.PermissionStateManager.ListAdmin()) == 0 {
			logger.Infof("您似乎正在部署全新的BOT，请通过qq对bot私聊发送<%v>(不含括号)获取管理员权限，然后私聊发送<%v>(不含括号)开始使用您的bot",
				l.CommandShowName(WhosyourdaddyCommand), l.CommandShowName(HelpCommand))
		}
	*/
}

func (b *BBot) Stop() {
	logger := log.CurrentModuleLogger()
	logger.Warn("bbot framework stopping ...")

	wg := sync.WaitGroup{}
	// stop Source modules
	for _, sourceModule := range b.Source {
		s := sourceModule
		go func() {
			wg.Add(1)
			_ = s.Stop()
			wg.Done()
		}()
	}
	// stop Platform modules
	for _, platformModule := range b.Platform {
		p := platformModule
		go func() {
			wg.Add(1)
			_ = p.Stop()
			wg.Done()
		}()
	}
	// stop Command modules
	// yet need to do nothing

	wg.Wait()
	logger.Info("stopped")
	b.Source = make([]Source, 0)
	b.Platform = make([]Platform, 0)
}

// GenerateExampleConfig gathers example configs from all plugins
// and write to bbot.yaml file
func (b *BBot) GenerateExampleConfig() (err error) {
	logger := log.CurrentModuleLogger()

	// TODO: gather example configs from all plugins
	//   and generate a config file
	err = os.WriteFile(consts.BBotConfigFilePath, []byte(utils.ExampleConfig()), 0755)
	if err != nil {
		logger.
			WithError(err).
			WithField("config_filename", consts.BBotConfigFilePath).
			Errorf("failed to generate example config")
		return err
	}
	logger.
		WithField("config_filename", consts.BBotConfigFilePath).
		Infof("最小配置已生成，请按需修改，如需高级配置请查看帮助文档")
	return err
}

// RegisterSource register a Source to BBot instance
func (b *BBot) RegisterSource(s Source) (err error) {
	s.NoticeTo(b.noticeChannel)
	b.Source = append(b.Source, s)
	return nil
}

// RegisterPlatform register a Platform to BBot instance
func (b *BBot) RegisterPlatform(p Platform) (err error) {
	b.Platform = append(b.Platform, p)
	return nil
}
