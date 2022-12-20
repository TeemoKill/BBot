package qq

import (
	"os"

	"github.com/TeemoKill/BBot/log"
	"github.com/TeemoKill/BBot/utils"

	miraiClient "github.com/Mrs4s/MiraiGo/client"
)

func GenRandomDevice() (err error) {
	logger := log.CurrentModuleLogger()

	miraiClient.GenRandomDevice()
	fileExist, err := utils.FileExist(DeviceInfoFilePath)
	if err != nil {
		logger.WithError(err).Errorf("check device info file error")
		return err
	}
	if fileExist {
		logger.WithField("deviceinfo_path", DeviceInfoFilePath).
			Warn("device info file exists, will not overwrite")
		return nil
	}

	err = os.WriteFile(DeviceInfoFilePath, miraiClient.SystemDeviceInfo.ToJson(), os.FileMode(0755))
	if err != nil {
		logger.WithError(err).Errorf("error during WriteFile")
		return err
	}

	return nil
}
