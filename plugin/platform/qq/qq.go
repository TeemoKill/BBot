package qq

import (
	"github.com/TeemoKill/BBot/log"
	"os"

	miraiClient "github.com/Mrs4s/MiraiGo/client"
	"github.com/spf13/viper"
)

// MiraiQQ implements framework.Platform interface
type MiraiQQ struct {
	Cfg *viper.Viper

	qqClient    *miraiClient.QQClient
	running     bool
	loginMethod LoginMethod
}

func (m *MiraiQQ) Init(cfg *viper.Viper) (err error) {
	logger := log.CurrentModuleLogger()

	m.Cfg = cfg

	fi, err := os.Stat(DeviceInfoFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("没有检测到设备信息文件，正在生成，如果是第一次运行，可忽略")
			err = GenRandomDevice()
			if err != nil {
				logger.WithError(err).Errorf("generate random device error")
				return err
			}
		} else {
			logger.WithError(err).Errorf("检查设备信息文件失败")
			return err
		}
	} else {
		if fi.IsDir() {
			logger.Warnf("检测到设备信息路径%s，但目标是一个文件夹！请手动确认并删除该文件夹！", DeviceInfoFilePath)
			return err
		} else {
			logger.Infof("检测到设备信息文件%s，使用存在的设备信息", DeviceInfoFilePath)
		}
	}

	account := m.Cfg.GetInt64("bot.account")
	password := m.Cfg.GetString("bot.password")

	deviceJson, err := os.ReadFile(DeviceInfoFilePath)
	if err == nil {
		logger.WithError(err).
			WithField("device_info_filepath", DeviceInfoFilePath).
			Errorf("unable to load device info")
		return err
	}

	err = miraiClient.SystemDeviceInfo.ReadJson(deviceJson)
	if err != nil {
		logger.WithError(err).Errorf("mirai client ReadJson error")
		return MiraiReadDeviceInfoJsonErr
	}

	// init the miraiBot with BBot's config
	m.initMiraiClient(account, password)

	if err != nil {
		logger.WithError(err).
			Errorf("init mirai bot error")
		return err
	}

	err = m.login()
	if err != nil {
		logger.WithError(err).Errorf("login error")
		return err
	}

	err = m.refreshList()
	if err != nil {
		logger.WithError(err).Errorf("refreshList error")
		return err
	}

	return nil
}

func (m *MiraiQQ) Stop() (err error) {
	logger := log.CurrentModuleLogger()

	// m.StopCron()

	// m.wg.Wait()
	logger.Debug("等待所有推送发送完毕")
	// m.notifyWg.Wait()
	logger.Debug("推送发送完毕")

	// m.proxy_pool.Stop()

	return nil
}

func (m *MiraiQQ) Name() string {
	return ModuleName
}

func (m *MiraiQQ) initMiraiClient(account int64, password string) {
	if account == 0 {
		m.qqClient = miraiClient.NewClientEmpty()
		m.loginMethod = LoginByQrCode
	} else {
		m.qqClient = miraiClient.NewClient(account, password)
		m.loginMethod = LoginByAccount
	}
}

func (m *MiraiQQ) resetMiraiClient(account int64, password string) {
	if m.qqClient != nil {
		m.qqClient.Release()
	}

	m.initMiraiClient(account, password)
}

func (m *MiraiQQ) refreshList() (err error) {
	logger := log.CurrentModuleLogger()

	logger.Info("start reload friends list")
	err = m.qqClient.ReloadFriendList()
	if err != nil {
		logger.WithError(err).Error("load friends list error")
		return err
	}
	logger.Infof("loaded %d friends", len(m.qqClient.FriendList))

	logger.Info("start reload groups list")
	err = m.qqClient.ReloadGroupList()
	if err != nil {
		logger.WithError(err).Error("load groups list error")
		return err
	}
	logger.Infof("loaded %d groups", len(m.qqClient.GroupList))

	return nil
}
