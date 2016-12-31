package mysql

// QueryBuilder 数据查询构造器
type QueryBuilder struct {
	Engine *DB
	table  string
	// 查询条件
	where  string
	args   []interface{}
	column []string
	// 排序条件
	order []string
}

// Table 数据查询的表
func (q *QueryBuilder) Table(table string) *QueryBuilder {
	q.table = table
	return q
}

// Where 数据查询条件
func (q *QueryBuilder) Where(sql string, args ...interface{}) *QueryBuilder {
	q.where = sql
	q.args = args
	return q
}

// OrderBy 数据排序条件
func (q *QueryBuilder) OrderBy(sql ...string) *QueryBuilder {
	q.order = sql
	return q
}

func (q *QueryBuilder) whereStr() string {
	if q.where != "" {
		return " WHERE " + q.where
	}
	return ""
}
