/*
Package mysql 数据库处理对象
*/
package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	pulic_type "github.com/asyou-me/lib.v1/pulic_type"
)

func init() {
	sql.Register("mysql_asyou", &mysql.MySQLDriver{})
}

// DB 数据库处理对象
type DB struct {
	*gorm.DB
	//数据操作连接池
	loger pulic_type.Logger
}

// Open 创建新的数据库对象
func (d *DB) Open(conf *pulic_type.MicroSerType) error {
	//初始化数据库
	var err error

	d.DB, err = gorm.Open("mysql_asyou", conf.Id+":"+
		conf.Secret+"@tcp("+conf.Addr+
		")/"+conf.Attr["Database"].(string)+"?charset=utf8&parseTime=True&loc=Local")
	// d.DB.LogMode(true)
	return err
}

// Stop 停止数据库
func (d *DB) Stop() error {
	return d.Close()
}

// NewDB 定义新建数据库操作对象的当法
func NewDB(conf *pulic_type.MicroSerType) (*DB, error) {
	db := DB{}
	err := db.Open(conf)
	return &db, err
}
