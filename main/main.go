package main

import (
	"github.com/TeemoKill/BBot/consts"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/TeemoKill/BBot/framework"
	"github.com/TeemoKill/BBot/log"
	"github.com/TeemoKill/BBot/storage"
	"github.com/TeemoKill/BBot/storage/buntdb"
	"github.com/TeemoKill/BBot/utils"
)

var sigChannel = make(chan os.Signal, 1)

func main() {
	var err error

	// -------- Init Logger --------
	logger, err := log.Init()
	if err != nil {
		return
	}

	// -------- Version Info --------
	logger.Infof("Tags: %v", utils.Tags)
	logger.Infof("Commit_ID: %v", utils.CommitId)
	logger.Infof("Build_Time: %v", utils.BuildTime)

	// -------- Init KV Storage --------
	err = buntdb.Init("")
	if err != nil {
		if err == storage.ErrLockNotHold {
			logger.Warnf("tryLock数据库失败：您可能重复启动了这个BOT！\n如果您确认没有重复启动，请尝试删除.lsp.db.lock文件并重新运行。")
		} else {
			logger.Warnf("无法正常初始化数据库！请检查.lsp.db文件权限是否正确，如无问题则为数据库文件损坏，请阅读文档获得帮助。")
		}
		return
	}
	if runtime.GOOS == "windows" {
		err = exitHook(func() { _ = buntdb.Close() })
		if err != nil {
			_ = buntdb.Close()
			logger.Warnf("无法正常初始化Windows环境！")
			return
		}
	} else {
		defer func() { _ = buntdb.Close() }()
	}

	// -------- Start BBot Service --------
	bBot := framework.New()
	err = bBot.Init()
	switch err {
	case consts.ExampleConfigGeneratedErr:
		logger.Infof("已生成最小配置，退出")
		return
	case nil:
		// no error
		break
	default:
		logger.WithError(err).Errorf("BBot 初始化失败")
		return
	}

	bBot.StartService()

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	<-sigChannel
	bBot.Stop()

}
