package system

import (
	"database/sql"
	"errors"
	"github.com/ejin66/goejin/util"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"sync"
)

var db *sql.DB

var dbLock sync.Mutex

type sqlBuilder struct {
	Sql string
}

type OrderBy struct {
	Name        string
	Orientation string
}

/*
Ipt : input type
key是数据库的column name , value是对应的值
*/
type Ipt map[string]interface{}

func Builder() *sqlBuilder {
	return new(sqlBuilder)
}

func valueFormat(value interface{}) (string, error) {
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
	default:
		return "", errors.New("can't recognize condition value type")
	}
	return result, nil
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

func concatWhereCondition(condition Ipt, conditionState ...string) (string, error) {
	var where string
	for key, value := range condition {
		result, err := valueFormat(value)
		if err != nil {
			return "", err
		}
		where += " AND `" + key + "` = " + result
	}

	for _, item := range conditionState {
		where += item
	}
	return where, nil
}

func (this *sqlBuilder) Select(table string, condition Ipt, conditionState ...string) *sqlBuilder {
	where, err := concatWhereCondition(condition, conditionState...)

	if err != nil {
		util.PrintError(err)
		return nil
	}

	sql := "SELECT * FROM " + table + " WHERE  1 = 1 " + where
	this.Sql = sql
	return this
}

func (this *sqlBuilder) OrderBy(orders ...OrderBy) *sqlBuilder {
	sql := this.Sql + " ORDER BY "
	for _, item := range orders {
		sql += item.Name + " " + item.Orientation + ","
	}
	this.Sql = sql[0 : len(sql)-1]
	return this
}

func (this *sqlBuilder) Limit(count int) *sqlBuilder {
	this.Sql += " LIMIT " + strconv.Itoa(count)
	return this
}

func (this *sqlBuilder) Limit2(start int, count int) *sqlBuilder {
	this.Sql += " LIMIT " + strconv.Itoa(start) + "," + strconv.Itoa(count)
	return this
}

func (this *sqlBuilder) SetSql(sql string) *sqlBuilder {
	this.Sql = sql
	return this
}

func (this *sqlBuilder) BuildSingle(model interface{}) error {
	util.Print(this.Sql)
	results := Query(this.Sql)
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
	util.Print(this.Sql)
	results := Query(this.Sql)
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

func Delete(table string, condition Ipt, conditionState ...string) bool {
	where, err := concatWhereCondition(condition, conditionState...)

	if err != nil {
		util.PrintError(err)
		return false
	}

	sql := "DELETE FROM " + table + " WHERE 1 = 1 " + where
	util.Print(sql)
	_, err2 := getDB().Exec(sql)

	if err2 != nil {
		util.PrintError(err2)
		return false
	}
	return true
}

func Update(table string, model interface{}, escapeColumns []string, condition Ipt, conditionState ...string) bool {
	where, err := concatWhereCondition(condition, conditionState...)

	if err != nil {
		util.PrintError(err)
		return false
	}

	var updateStatement string
	elements := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()
	for i := 0; i < elements.NumField(); i++ {
		key := elements.Field(i).Tag.Get("db")
		value := modelValue.Field(i).Interface()
		result, err := valueFormat(value)
		if err != nil {
			continue
		}
		if util.InArray(escapeColumns, key) {
			continue
		}
		updateStatement += ", `" + key + "` = " + result
	}

	updateStatement = string(updateStatement[1:])

	sql := "UPDATE " + table + " SET " + updateStatement + " WHERE 1= 1 " + where

	util.Print(sql)

	_, err3 := getDB().Exec(sql)

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
		result, err := valueFormat(value)
		if err != nil {
			continue
		}
		if util.InArray(escapeColumns, key) {
			continue
		}
		columns += ",`" + key + "`"
		values += "," + result
	}

	columns = string(columns[1:])
	values = string(values[1:])

	sql := "INSERT INTO " + table + "(" + columns + ") VALUES" + "(" + values + ")"
	util.Print(sql)

	result, err2 := getDB().Exec(sql)

	if err2 != nil {
		util.PrintError(err2)
		return -1
	}
	i, _ := result.LastInsertId()
	return i

}

func connect() {
	dataSourceName := GetConfig().DbUser + ":" + GetConfig().DbPassword + "@tcp(" + GetConfig().DbAddress + ":" + GetConfig().DbPort + ")/" + GetConfig().DbName + "?charset=utf8"
	//util.Print(dataSourceName)
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
