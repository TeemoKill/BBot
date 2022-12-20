package buntdb

import (
	"fmt"

	"github.com/TeemoKill/BBot/storage"

	"github.com/gofrs/flock"
	"github.com/modern-go/gls"
	"github.com/tidwall/buntdb"
)

var instance *buntdb.DB
var fileLock *flock.Flock

// Init 初始化buntdb，正常情况下框架会负责初始化
func Init(dbPath string) error {
	if dbPath == "" {
		dbPath = StoragePathDefault
	}
	if dbPath != StoragePathMemory {
		var dbLock = dbPath + ".lock"
		fileLock = flock.New(dbLock)
		ok, err := fileLock.TryLock()
		if err != nil {
			fmt.Printf("buntdb tryLock err: %v", err)
		}
		if !ok {
			return storage.ErrLockNotHold
		}
	}
	buntDB, err := buntdb.Open(dbPath)
	if err != nil {
		return err
	}
	if dbPath != StoragePathMemory {
		buntDB.SetConfig(buntdb.Config{
			SyncPolicy:           buntdb.EverySecond,
			AutoShrinkPercentage: 10,
			AutoShrinkMinSize:    1 * 1024 * 1024,
		})
	}
	instance = buntDB
	return nil
}

// GetClient 获取 buntdb.DB 对象，如果没有初始化会返回 ErrNotInitialized
func GetClient() (*buntdb.DB, error) {
	if instance == nil {
		return nil, ErrNotInitialized
	}
	return instance, nil
}

// MustGetClient 获取 buntdb.DB 对象，如果没有初始化会panic，在编写订阅组件时可以放心调用
func MustGetClient() *buntdb.DB {
	if instance == nil {
		panic(ErrNotInitialized)
	}
	return instance
}

// Close 关闭buntdb，正常情况下框架会负责关闭
func Close() error {
	if instance != nil {
		if itx := gls.Get(txKey); itx != nil {
			itx.(*buntdb.Tx).Rollback()
		}
		err := instance.Close()
		if err != nil {
			return err
		}
		instance = nil
	}
	if fileLock != nil {
		return fileLock.Unlock()
	}
	return nil
}
