package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// 类型
const (
	Invalid       uint8 = iota // 0
	Bool                       // 1
	Int                        // 2
	Int8                       // 3
	Int16                      // 4
	Int32                      // 5
	Int64                      // 6
	Uint                       // 7
	Uint8                      // 8
	Uint16                     // 9
	Uint32                     // 10
	Uint64                     // 11
	Uintptr                    // 12
	Float32                    // 13
	Float64                    // 14
	Complex64                  // 15
	Complex128                 // 16
	Array                      // 17
	Chan                       // 18
	Func                       // 19
	Interface                  // 20
	Map                        // 21
	Ptr                        // 22
	Slice                      // 23
	String                     // 24
	Struct                     // 25
	UnsafePointer              // 26
	IntArray                   // 27
	StringArray                // 28
)

// V 基础类型
type V struct {
	T           uint8
	V           string
	IntArray    *[]int64
	StringArray *[]string
}

// MarshalJSON 序列化时调用
func (v *V) MarshalJSON() ([]byte, error) {
	switch v.T {
	case Int, Int8, Int16, Int32, Int64:
		if len(v.V) == 0 {
			return []byte{'0'}, nil
		}
		return []byte(v.V), nil
	case Bool:
		if v.V == "true" {
			return []byte(v.V), nil
		}
		return []byte("false"), nil
	case String:
		return []byte(`"` + v.V + `"`), nil
	case IntArray:
		if v.IntArray == nil {
			return []byte("[]"), nil
		}
		datas := *v.IntArray
		if len(datas) == 0 {
			return []byte("[]"), nil
		}
		if len(datas) == 1 {
			return []byte("[" + fmt.Sprint(datas[0]) + "]"), nil
		}

		strs := make([]string, len(datas))
		n := len(datas) - 1
		for k, v := range datas {
			str := fmt.Sprint(v)
			strs[k] = str
			n = n + len(str)
		}

		b := make([]byte, n+2)
		bp := 0
		bp += copy(b, "[")
		bp += copy(b[bp:], strs[0])
		for _, v := range strs[1:] {
			bp += copy(b[bp:], ",")
			bp += copy(b[bp:], v)
		}
		bp += copy(b[bp:], "]")
		return b, nil
	case StringArray:
		if v.StringArray == nil {
			return []byte("[]"), nil
		}
		bt, _ := json.Marshal(v.StringArray)
		return bt, nil
	default:
		return []byte{}, errors.New("无法识别类型:" + fmt.Sprint(v.T))
	}
}

// Scan 渲染数据到结构
func (v *V) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		if v.T != String {
			return errors.New("结构体字段为 string")
		}
		v.V = s
		return nil
	case time.Time:
		if v.T != Int64 {
			return errors.New("结构体字段为 string")
		}
		v.V = fmt.Sprint(s.Unix())
		return nil
	case nil:
		return nil
	}

	switch v.T {
	case String:
		v.V = asString(value)
		return nil
	case IntArray:
		by := []byte(asString(value))
		v.IntArray = new([]int64)
		json.Unmarshal(by, v.IntArray)
		return nil
	case StringArray:
		by := []byte(asString(value))
		v.StringArray = new([]string)
		json.Unmarshal(by, v.StringArray)
		return nil
	case Int16, Int8, Int64, Int32, Uint16, Uint8, Uint64, Uint32, Int:
		v.V = asString(value)
		return nil
	case Bool:
		bv, err := driver.Bool.ConvertValue(value)
		if err == nil {
			v.V = fmt.Sprint(bv.(bool))
		}
		return err
	}

	return errors.New("无法解析类型")
}

// Value 写入结构数据到数据库
func (v V) Value() (driver.Value, error) {
	if v.T == IntArray {
		if v.IntArray == nil {
			v.IntArray = &[]int64{}
		}
		bt, _ := json.Marshal(v.IntArray)
		return string(bt), nil
	}
	if v.T == StringArray {
		if v.StringArray == nil {
			v.StringArray = &[]string{}
		}
		bt, _ := json.Marshal(v.StringArray)
		return string(bt), nil
	}
	if v.T == Bool {
		if v.V == "true" {
			return true, nil
		}
		return false, nil
	}
	return v.V, nil
}

// InsertV 插入数据到数据库
func (q *QueryBuilder) InsertV(data map[string]*V) (string, error) {
	var inserts = ""
	var fields = ""
	var lenData = len(data)
	var indexData = lenData - 1
	var values = make([]interface{}, lenData)
	var index int
	for k, v := range data {
		if index == indexData {
			fields = fields + "`" + k + "`"
			inserts = inserts + `?`
		} else {
			fields = fields + "`" + k + "`,"
			inserts = inserts + `?` + `,`
		}
		values[index] = v
		index = index + 1
	}
	db := q.Engine.Exec(`INSERT INTO `+q.table+` (`+fields+`) VALUES (`+inserts+`)`, values...)
	return "", db.Error
}

// UpdateV 更新数据到数据库
func (q *QueryBuilder) UpdateV(data map[string]*V) (string, error) {
	if q.where == "" {
		return "", errors.New("更新条件不能为空")
	}
	var sets = ""
	var lenData = len(data)
	var indexData = lenData - 1
	var index int
	var values = make([]interface{}, lenData)
	for k, v := range data {
		if index == indexData {
			sets = sets + "`" + k + "`" + `=?`
		} else {
			sets = sets + "`" + k + "`" + `=?,`
		}
		values[index] = v
		index = index + 1
	}
	db := q.Engine.Exec(`UPDATE `+q.table+` SET `+sets+q.whereStr(), values...)
	if db.Error != nil {
		return "0", db.Error
	}
	return fmt.Sprint(db.RowsAffected), db.Error
}

// GetV 获取字段数据
func (q *QueryBuilder) GetV(data map[string]*V) (err error) {
	var gets = ""
	var values = make([]interface{}, len(data))
	var indexData = len(data) - 1
	var index = 0
	for k, v := range data {
		if index == indexData {
			gets = gets + "`" + k + "`"
		} else {
			gets = gets + "`" + k + "`,"
		}
		values[index] = v
		index = index + 1
	}
	query := `SELECT ` + gets + ` FROM ` + q.table + q.whereStr() + q.OrderStr()
	row := q.Engine.Raw(query).Row()
	err = row.Scan(values...)
	if err != nil {
		return
	}
	return nil
}

// ListV 获取字段数据
func (q *QueryBuilder) ListV(files map[string]uint8, limit int, offset int) ([]*map[string]*V, error) {
	var gets string
	var indexData = len(files) - 1
	var index = 0
	var fs = make([]string, len(files))
	for k := range files {
		if index == indexData {
			gets = gets + "`" + k + "`"
		} else {
			gets = gets + "`" + k + "`,"
		}
		fs[index] = k
		index = index + 1
	}
	outs := make([]*map[string]*V, 0, limit)
	query := `SELECT ` + gets + ` FROM ` + q.table + q.whereStr() + q.OrderStr() + ` LIMIT ` + fmt.Sprint(limit) + ` OFFSET ` + fmt.Sprint(offset)
	db := q.Engine.Raw(query)
	if db.Error != nil {
		return nil, db.Error
	}
	rows, err := db.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := map[string]*V{}
		values := make([]interface{}, len(files))
		i := 0
		for _, k := range fs {
			v := &V{
				T: files[k],
			}
			item[k] = v
			values[i] = v
			i = i + 1
		}
		err = rows.Scan(values...)
		if err != nil {
			continue
		}
		outs = append(outs, &item)
	}
	return outs, nil
}

// Search 获取字段数据
func (q *QueryBuilder) Search(files map[string]uint8) ([]*map[string]*V, error) {
	var gets string
	var indexData = len(files) - 1
	var index = 0
	var fs = make([]string, len(files))
	for k := range files {
		if index == indexData {
			gets = gets + "`" + k + "`"
		} else {
			gets = gets + "`" + k + "`,"
		}
		fs[index] = k
		index = index + 1
	}
	outs := make([]*map[string]*V, 0, 10)
	query := `SELECT ` + gets + ` FROM ` + q.table + q.whereStr() + q.OrderStr()
	db := q.Engine.Raw(query)
	if db.Error != nil {
		return nil, db.Error
	}
	rows, err := db.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := map[string]*V{}
		values := make([]interface{}, len(files))
		i := 0
		for _, k := range fs {
			v := &V{
				T: files[k],
			}
			item[k] = v
			values[i] = v
			i = i + 1
		}
		err = rows.Scan(values...)
		if err != nil {
			continue
		}
		outs = append(outs, &item)
	}
	return outs, nil
}
