package database

import (
	"environment/cfgargs"
)

// DB 数据库+缓存
type DB struct {
}

// NewDB create
func NewDB(cfg *cfgargs.SrvConfig) (*DB, error) {
	db := &DB{}
	return db, db.init(cfg)
}

func (d *DB) init(cfg *cfgargs.SrvConfig) error {
	return nil
}
