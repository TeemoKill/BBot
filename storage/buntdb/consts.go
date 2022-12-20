package buntdb

import "errors"

const StoragePathDefault = ".bbot.bunt.db"
const StoragePathMemory = ":memory:"

// -------- Errors --------

var (
	ErrKeyExist       = errors.New("key exist")
	ErrNotInitialized = errors.New("not initialized")
	ErrRollback       = errors.New("rollback")
)
