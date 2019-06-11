package system

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ejin66/goejin/util"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"sync"
)

var db *sql.DB

var dbLock sync.Mutex

type sqlBuilder struct {
	sql  string
	args []interface{}
}

type OrderBy struct {
	Column
	Direct
}

type Condition struct {
	Logic
	Sign
	Column
	Value interface{}
	next  *Condition
	last  *Condition
}

type Column struct {
	Table string
	Key   string
}

type Logic int

type Sign int

type Join int

type Direct int

const (
	LogicNone Logic = iota
	LogicOr
	LogicAnd
)

const (
	SignEqual Sign = iota
	SignGreater
	SignLess
)

const (
	InnerJoin Join = iota
	LeftJoin
	RightJoin
)

const (
	ASC Direct = iota
	DESC
)

func Builder() *sqlBuilder {
	return new(sqlBuilder)
}

func valueFormat(value interface{}) string {
	var result string
	switch value.(type) {
	case int:
		result = strconv.Itoa(value.(int))
	case string:
		result = "'" + value.(string) + "'"
	case bool:
		if value.(bool) {
			result = "1"
		} else {
			result = "0"
		}
	case Column:
		c := value.(Column)
		result = c.parse()
	default:
		result = "'" + value.(string) + "'"
	}
	return result
}

func valueFormat2(typeName string, value string) (result interface{}) {
	switch typeName {
	case "int":
		result, _ = strconv.Atoi(value)
	case "bool":
		if value == "1" {
			result = true
		} else {
			result = false
		}
	default:
		result = value
	}
	return
}

func concatWhereCondition(conditions ...Condition) string {
	var where string
	for _, v := range conditions {
		where += v.parseCondition()
	}

	return where
}

func (c Condition) Append(conditions ...Condition) *Condition {
	var current = &c
	for _, item := range conditions {
		current.next = &item
		item.last = current
		current = current.next
	}
	return &c
}

func (c *Condition) parseCondition() string {
	v := valueFormat(c.Value)

	var logic string

	switch c.Logic {
	case LogicOr:
		logic = "OR"
	case LogicAnd:
		logic = "AND"
	case LogicNone:
	default:
		logic = ""
	}

	var sign string
	switch c.Sign {
	case SignEqual:
		sign = "="
	case SignGreater:
		sign = ">"
	case SignLess:
		sign = "<"
	}

	var bracketLeft string
	var bracketRight string
	if c.last == nil {
		bracketLeft = "("
		bracketRight = ")"
	}

	var nextCondition string
	if c.next != nil {
		nextCondition = c.next.parseCondition()
	}

	result := " " + logic + " " + bracketLeft + c.Column.parse() + sign + v + nextCondition + bracketRight + " "
	return result
}

func (c *Column) parse() string {
	if len(c.Table) == 0 {
		return "`" + c.Key + "`"
	} else {
		return c.Table + ".`" + c.Key + "`"
	}
}

func (this *sqlBuilder) Select(table string, columns ...Column) *sqlBuilder {
	var columnStr string

	for _, v := range columns {
		columnStr += v.parse() + ","
	}

	if len(columnStr) == 0 {
		columnStr = "*"
	} else {
		columnStr = columnStr[0 : len(columnStr)-1]
	}

	this.sql = fmt.Sprintf("SELECT %s FROM %s ", columnStr, table)
	return this
}

func (this *sqlBuilder) Join(join Join, Table string, conditions ...Condition) *sqlBuilder {
	var joinStr string

	switch join {
	case InnerJoin:
		joinStr = "INNER JOIN"
	case LeftJoin:
		joinStr = "LEFT JOIN"
	case RightJoin:
		joinStr = "RIGHT JOIN"
	}

	this.sql += " " + joinStr + " " + Table + " ON " + concatWhereCondition(conditions...)
	return this
}

func (this *sqlBuilder) Where(conditions ...Condition) *sqlBuilder {
	this.sql += " WHERE " + concatWhereCondition(conditions...)
	return this
}

func (this *sqlBuilder) OrderBy(orders ...OrderBy) *sqlBuilder {
	query := this.sql + " ORDER BY "

	for _, item := range orders {
		var direct string

		if item.Direct == ASC {
			direct = "ASC"
		} else {
			direct = "DESC"
		}

		query += item.Column.parse() + " " + direct + ","
	}
	this.sql = query[0 : len(query)-1]
	return this
}

func (this *sqlBuilder) Limit(count int) *sqlBuilder {
	this.sql += " LIMIT " + strconv.Itoa(count)
	return this
}

func (this *sqlBuilder) LimitIndex(start int, count int) *sqlBuilder {
	this.sql += " LIMIT " + strconv.Itoa(start) + "," + strconv.Itoa(count)
	return this
}

func (this *sqlBuilder) SetSql(query string, values ...interface{}) *sqlBuilder {
	this.sql = query
	this.args = values
	return this
}

func (this *sqlBuilder) BuildExec() (sql.Result, error) {
	util.Print(this.sql)
	return getDB().Exec(this.sql, this.args...)
}

func (this *sqlBuilder) BuildSingle(model interface{}) error {
	util.Print(this.sql)
	results := Query(this.sql)
	if len(results) != 1 {
		return errors.New("query result size is not single!")
	}
	modelType := reflect.TypeOf(model)
	elements := modelType.Elem()

	for j := 0; j < elements.NumField(); j++ {
		field := elements.Field(j)
		if field.Tag.Get("db") == "" {
			continue
		}
		if _, ok := results[0][field.Tag.Get("db")]; !ok {
			continue
		}
		reflect.ValueOf(model).Elem().FieldByName(field.Name).Set(reflect.ValueOf(valueFormat2(field.Type.Name(),
			results[0][field.Tag.Get("db")])))
	}
	return nil
}

/*
 model's type must be *slice
*/
func (this *sqlBuilder) Build(model interface{}) {
	util.Print(this.sql)
	results := Query(this.sql)
	elementType := reflect.TypeOf(model).Elem().Elem()
	sliceV := reflect.ValueOf(model).Elem()

	for _, result := range results {
		modelTemp := reflect.New(elementType)
		for j := 0; j < elementType.NumField(); j++ {
			field := elementType.Field(j)
			if field.Tag.Get("db") == "" {
				continue
			}
			if _, ok := result[field.Tag.Get("db")]; !ok {
				continue
			}
			modelTemp.Elem().FieldByName(field.Name).Set(reflect.ValueOf(valueFormat2(field.Type.Name(), result[field.Tag.Get("db")])))
		}
		sliceV = reflect.Append(sliceV, modelTemp.Elem())
	}
	reflect.ValueOf(model).Elem().Set(sliceV)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return getDB().Exec(query, args...)
}

func Query(query string, args ...interface{}) []map[string]string {
	rows, err := getDB().Query(query, args...)

	if err != nil {
		util.PrintError(err)
		return []map[string]string{}
	}

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	var results []map[string]string

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, v := range values {
			if v != nil {
				record[columns[i]] = string(v.([]byte))
			}
		}
		results = append(results, record)
	}
	return results
}

func Delete(table string, conditions ...Condition) bool {
	where := concatWhereCondition(conditions...)

	query := "DELETE FROM " + table + " WHERE " + where
	util.Print(query)
	_, err2 := getDB().Exec(query)

	if err2 != nil {
		util.PrintError(err2)
		return false
	}
	return true
}

func Update(table string, model interface{}, escapeColumns []string, conditions ...Condition) bool {
	where := concatWhereCondition(conditions...)

	var updateStatement string
	elements := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()
	for i := 0; i < elements.NumField(); i++ {
		key := elements.Field(i).Tag.Get("db")
		value := modelValue.Field(i).Interface()
		result := valueFormat(value)

		if util.InArray(escapeColumns, key) {
			continue
		}
		updateStatement += ", `" + key + "` = " + result
	}

	updateStatement = string(updateStatement[1:])

	query := "UPDATE " + table + " SET " + updateStatement + " WHERE " + where

	util.Print(query)

	_, err3 := getDB().Exec(query)

	if err3 != nil {
		util.PrintError(err3)
		return false
	}
	return true
}

func Insert(table string, model interface{}, escapeColumns []string) int64 {
	var columns string
	var values string
	elements := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()
	for i := 0; i < elements.NumField(); i++ {
		key := elements.Field(i).Tag.Get("db")
		value := modelValue.Field(i).Interface()
		result := valueFormat(value)

		if util.InArray(escapeColumns, key) {
			continue
		}
		columns += ",`" + key + "`"
		values += "," + result
	}

	columns = string(columns[1:])
	values = string(values[1:])

	query := "INSERT INTO " + table + "(" + columns + ") VALUES" + "(" + values + ")"
	util.Print(query)

	result, err2 := getDB().Exec(query)

	if err2 != nil {
		util.PrintError(err2)
		return -1
	}
	i, _ := result.LastInsertId()
	return i
}

func connect() {
	dataSourceName := GetConfig().DbUser + ":" + GetConfig().DbPassword + "@tcp(" + GetConfig().DbAddress + ":" + GetConfig().DbPort + ")/" + GetConfig().DbName + "?charset=utf8"
	database, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		util.PrintError(err)
	} else {
		db = database
	}
}

func getDB() *sql.DB {
	if db == nil {
		dbLock.Lock()
		defer dbLock.Unlock()
		if db == nil {
			connect()
		}
	}
	return db
}
