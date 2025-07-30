package boltdb

import "go.etcd.io/bbolt"

// DB returns the underlying bbolt.DB instance. This is intended for testing purposes only.
func (bdl *BoltDBDataLayer) GetDb_ForTest() *bbolt.DB {
	return bdl.db
}
