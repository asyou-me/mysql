package mysql

// QueryBuilder 数据查询构造器
type QueryBuilder struct {
	Engine *DB
	table  string
	// 查询条件
	where string
	// 排序条件
	order  string
	column []string
}

// Table 数据查询的表
func (q *QueryBuilder) Table(table string) *QueryBuilder {
	q.table = table
	return q
}

// Where 数据查询条件
func (q *QueryBuilder) Where(sql string, args ...interface{}) *QueryBuilder {
	q.where = sql
	return q
}

// Order 排序条件
func (q *QueryBuilder) Order(sql string) *QueryBuilder {
	q.order = sql
	return q
}

// OrderStr 数据排序条件
func (q *QueryBuilder) OrderStr() string {
	if q.order != "" {
		return " ORDER BY " + q.order
	}
	return ""
}

func (q *QueryBuilder) whereStr() string {
	if q.where != "" {
		return " WHERE " + q.where
	}
	return ""
}
