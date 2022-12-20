package buntdb

import (
	"github.com/tidwall/buntdb"
)

func IsRollback(e error) bool {
	return e == ErrRollback
}

func IsNotFound(e error) bool {
	return e == buntdb.ErrNotFound
}
