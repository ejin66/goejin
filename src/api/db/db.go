package db

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"api/config"
	"strconv"
	"api/system"
	"errors"
	"fmt"
)

var db *sql.DB

/*
Ipt : input type
key是数据库的column name , value是对应的值
 */
type Ipt map[string]interface{}

func init() {
	dataSourceName := config.DB_USER + ":" + config.DB_PASSWORD + "@tcp(localhost:" + config.DB_PORT + ")/" + config.DB_NAME + "?charset=utf8"
	//fmt.Println(dataSourceName)
	database, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		system.PrintError(err)
	} else {
		db = database
	}
}

func valueFormat(value interface{}) (string, error) {
	var result string
	switch value.(type) {
	case int:
		result = strconv.Itoa(value.(int))
	case string:
		result = "'" + value.(string) + "'"
	default:
		return "", errors.New("can't recognize condition value type")
	}
	return result, nil
}

func concatWhereCondition(condition Ipt, conditionState ...string) (string, error) {
	var where string
	for key, value := range condition {
		result, err := valueFormat(value)
		if err != nil {
			return "", err
		}
		where += " AND " + key + " = " + result
	}

	for _, item := range conditionState {
		where += item
	}
	return where, nil
}

func Query(table string, condition Ipt, conditionState ...string) []map[string]string {

	where, err := concatWhereCondition(condition, conditionState...)

	if err != nil {
		system.PrintError(err)
		return nil
	}

	sql := "SELECT * FROM " + table + " WHERE  1 = 1 " + where

	//fmt.Println(sql)

	rows, err := db.Query(sql)

	if err != nil {
		system.PrintError(err)
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
		system.PrintError(err)
		return false
	}

	sql := "DELETE FROM " + table + " WHERE 1 = 1 " + where
	fmt.Println(sql)
	_, err2 := db.Exec(sql)

	if err2 != nil {
		system.PrintError(err2)
		return false
	}
	return true
}

func Update(table string, update Ipt, condition Ipt, conditionState ...string) bool {
	where, err := concatWhereCondition(condition, conditionState...)

	if err != nil {
		system.PrintError(err)
		return false
	}

	var updateStatement string
	for k, v := range update {
		result, err2 := valueFormat(v)

		if err2 != nil {
			system.PrintError(err2)
			return false
		}

		updateStatement += ", " + k + " = " + result
	}

	updateStatement = string(updateStatement[1:])

	sql := "UPDATE " + table + " SET " + updateStatement + " WHERE 1= 1 " + where

	fmt.Println(sql)

	_, err3 := db.Exec(sql)

	if err3 != nil {
		system.PrintError(err3)
		return false
	}
	return true
}

func Insert(table string, insert Ipt) int64 {

	var columns string
	var values string

	for k, v := range insert {
		result, err := valueFormat(v)
		if err != nil {
			system.PrintError(err)
			return -1
		}
		columns += "," + k
		values += "," + result
	}

	columns = string(columns[1:])
	values = string(values[1:])

	sql := "INSERT INTO " + table + "(" + columns + ") VALUES" + "("  + values + ")"
	fmt.Println(sql)

	result, err2 := db.Exec(sql)

	if err2 != nil {
		system.PrintError(err2)
		return -1
	}
	i,_ := result.LastInsertId()
	return i

}
