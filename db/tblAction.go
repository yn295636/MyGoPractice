package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func InsertInto(dbConn *sql.DB, tblName string, kvMap map[string]interface{}) (sql.Result, error) {
	var fields []string
	var values []interface{}
	var valuePlaceholders []string
	for k, v := range kvMap {
		fields = append(fields, k)
		values = append(values, v)
		valuePlaceholders = append(valuePlaceholders, "?")
	}
	sqlStat := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
		tblName, strings.Join(fields, ","), strings.Join(valuePlaceholders, ","))
	return dbConn.Exec(sqlStat, values...)
}

func InsertIntoByTx(tx *sql.Tx, tblName string, kvMap map[string]interface{}) (sql.Result, error) {
	var fields []string
	var values []interface{}
	var valuePlaceholders []string
	for k, v := range kvMap {
		fields = append(fields, k)
		values = append(values, v)
		valuePlaceholders = append(valuePlaceholders, "?")
	}
	sqlStat := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
		tblName, strings.Join(fields, ","), strings.Join(valuePlaceholders, ","))
	return tx.Exec(sqlStat, values...)
}

func alterTableByTx(dbTx *sql.Tx, tableName, modifiedSchema string) error {
	sqlStat := fmt.Sprintf("ALTER TABLE %v \n%v",
		tableName, modifiedSchema)
	_, err := dbTx.Exec(sqlStat)
	return err
}

func createTable(dbConn *sql.DB, tableName, tableSchema string) error {
	sqlStat := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v(\n%v\n)\nCOLLATE='utf8mb4_unicode_ci'\nENGINE=InnoDB",
		tableName, tableSchema)
	_, err := dbConn.Exec(sqlStat)
	if err != nil {
		log.Printf("Create table %v failed, %v", tableName, err)
		return err
	}
	return nil
}

func createTableByTx(dbTx *sql.Tx, tableName, tableSchema string) error {
	sqlStat := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v(\n%v\n)\nCOLLATE='utf8mb4_unicode_ci'\nENGINE=InnoDB",
		tableName, tableSchema)
	_, err := dbTx.Exec(sqlStat)
	return err
}
