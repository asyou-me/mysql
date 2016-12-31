package mysql

import "fmt"

// Del 删除数据
//
// table:删除数据的表
//
// req:条件sql写法 where xxx
func (d *DB) Del(table string, req string) (err error) {
	if req != "" {
		req = "WHERE " + req
	}
	sql := `DELETE FROM ` + table + ` ` + req
	db := d.Exec(sql)
	if db.Error != nil {
		return
	}
	return
}

// Count 获取数据的条数
//
// table:数据的表
//
// req:条件sql写法 where xxx
func (d *DB) Count(table string, where string) int64 {
	var re int64
	fmt.Println(`SELECT COUNT(*) FROM ` + table + ` ` + where)
	db := d.DB.Table(table).Count(&re)
	if db.Error != nil {
		return 0
	}
	return re
}

// Table 建立一个 针对于 table 表的数据库查询对象
//
// table:数据库表名
func (d *DB) Table(table string) *QueryBuilder {
	builder := &QueryBuilder{
		Engine: d,
	}
	builder.Table(table)
	return builder
}
