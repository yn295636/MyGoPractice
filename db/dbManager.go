package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
)

type TDbVersion int

const (
	DbMyGoPractice = "mygopractice"
	tblMeta        = "meta"
	TblUser        = "user"
	TblPhone       = "phone"
	metaSchema     = "version INT(8) NOT NULL UNIQUE, update_time INT(10) UNSIGNED NOT NULL"

	DbVer1      TDbVersion = 1
	DbVer2      TDbVersion = 2
	DbVer3      TDbVersion = 3
	DbLatestVer TDbVersion = DbVer3
)

var mapDB sync.Map

func connectDbRoot(dbAddr, dbUser, dbPassword string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%v:%v@tcp(%v)/?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		dbUser, dbPassword, dbAddr))
	if err != nil {
		log.Printf("Connect to mysql error, %v", err)
		return nil, err
	}
	return db, nil
}

func createDb(dbConn *sql.DB, dbName string) error {
	sqlStat := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", dbName)
	_, err := dbConn.Exec(sqlStat)
	if err != nil {
		log.Printf("Create db %v failed, %v", dbName, err)
		return err
	}
	return nil
}

func connectDbWithName(dbAddr, dbUser, dbPassword, dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		dbUser, dbPassword, dbAddr, dbName))
	if err != nil {
		log.Printf("Connect to mysql error, %v", err)
		return nil, err
	}
	return db, nil
}

func getCurrentDbVersion(dbConn *sql.DB) (TDbVersion, error) {
	var result TDbVersion
	sqlState := fmt.Sprintf("SELECT version FROM %v ORDER BY version DESC LIMIT 1", tblMeta)
	row := dbConn.QueryRow(sqlState)
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		log.Printf("Query db version failed, %v", err)
		return 1000000, err
	}
	return result, nil
}

func updateDbVersion(dbTx *sql.Tx, ver TDbVersion) error {
	kvMap := map[string]interface{}{
		"version":     ver,
		"update_time": time.Now().Unix(),
	}
	_, err := InsertIntoByTx(dbTx, tblMeta, kvMap)
	return err
}

func updateMyGoPracticeToVer1(dbConn *sql.DB) error {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Printf("DbMyGoPractice Tx begin error, %v", err)
		return err
	}
	defer TxPost(tx, err)

	userSchema := `id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			username VARCHAR(16) NOT NULL,
			PRIMARY KEY (id)`
	err = createTableByTx(tx, TblUser, userSchema)
	if err != nil {
		log.Printf("Create table user failed, %v", err)
		return err
	}

	err = updateDbVersion(tx, DbVer1)
	if err != nil {
		log.Printf("Update db version failed, %v", err)
		return err
	}
	return nil
}

func updateMyGoPracticeToVer2(dbConn *sql.DB) error {
	var err error
	tx, err := dbConn.Begin()
	if err != nil {
		log.Printf("DbMyGoPractice Tx begin error, %v", err)
		return err
	}
	defer TxPost(tx, err)

	phoneSchema := `id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			phone VARCHAR(16) NOT NULL,
			uid BIGINT(20) UNSIGNED NOT NULL,
			PRIMARY KEY (id)`
	err = createTableByTx(tx, TblPhone, phoneSchema)
	if err != nil {
		log.Printf("Create table phone failed, %v", err)
		return err
	}

	userAddSchema := `ADD COLUMN gender TINYINT(3) NOT NULL DEFAULT 0 COMMENT '0:unknown, 1:male, 2:female',
			ADD COLUMN age TINYINT(3) NOT NULL DEFAULT 0`
	err = alterTableByTx(tx, TblUser, userAddSchema)
	if err != nil {
		log.Printf("Alter table user failed, %v", err)
		return err
	}

	err = updateDbVersion(tx, DbVer2)
	if err != nil {
		log.Printf("Update db version failed, %v", err)
		return err
	}
	return nil
}

func updateMyGoPracticeToVer3(dbConn *sql.DB) error {
	var err error
	tx, err := dbConn.Begin()
	if err != nil {
		log.Printf("DbMyGoPractice Tx begin error, %v", err)
		return err
	}
	defer TxPost(tx, err)

	userAlterSchema := `ADD CONSTRAINT UNIQUE INDEX idx_username (username)`
	err = alterTableByTx(tx, TblUser, userAlterSchema)
	if err != nil {
		log.Printf("Alter table user failed, %v", err)
		return err
	}

	phoneAlterSchema := `ADD CONSTRAINT UNIQUE INDEX idx_phone (phone)`
	err = alterTableByTx(tx, TblPhone, phoneAlterSchema)
	if err != nil {
		log.Printf("Alter table phone failed, %v", err)
		return err
	}

	err = updateDbVersion(tx, DbVer3)
	if err != nil {
		log.Printf("Update db version failed, %v", err)
		return err
	}
	return nil
}

func createAllDbs(dbAddr, dbUser, dbPassword string) error {
	db, err := connectDbRoot(dbAddr, dbUser, dbPassword)
	if err != nil {
		log.Printf("Connect db root path failed, %v", err)
		return err
	}
	defer db.Close()

	err = createDb(db, DbMyGoPractice)
	if err != nil {
		log.Printf("Create DbMyGoPractice failed, %v", err)
		return err
	}
	return nil
}

func InitDb(dbAddr, dbUser, dbPassword string, targetVersion TDbVersion) error {
	err := createAllDbs(dbAddr, dbUser, dbPassword)
	if err != nil {
		log.Printf("Create db failed, %v", err)
		return err
	}

	dbMyGoPractice, err := connectDbWithName(dbAddr, dbUser, dbPassword, DbMyGoPractice)
	if err != nil {
		log.Printf("Connect db DbMyGoPractice failed, %v", err)
		return err
	}

	err = createTable(dbMyGoPractice, tblMeta, metaSchema)
	if err != nil {
		log.Printf("Create table meta failed, %v", err)
		return err
	}

	curVersion, err := getCurrentDbVersion(dbMyGoPractice)
	if err != nil {
		log.Printf("Get current version of DbMyGoPractice failed, %v", err)
		return err
	}

	if curVersion < DbVer1 && DbVer1 <= targetVersion {
		err = updateMyGoPracticeToVer1(dbMyGoPractice)
		if err != nil {
			log.Printf("Update DbMyGoPractice to version 1 failed, %v", err)
			return err
		}
	}

	if curVersion < DbVer2 && DbVer2 <= targetVersion {
		err = updateMyGoPracticeToVer2(dbMyGoPractice)
		if err != nil {
			log.Printf("Update DbMyGoPractice to version 2 failed, %v", err)
			return err
		}
	}

	if curVersion < DbVer3 && DbVer3 <= targetVersion {
		err = updateMyGoPracticeToVer3(dbMyGoPractice)
		if err != nil {
			log.Printf("Update DbMyGoPractice to version 3 failed, %v", err)
			return err
		}
	}

	mapDB.Store(DbMyGoPractice, dbMyGoPractice)
	return nil
}

func TxPost(tx *sql.Tx, err error) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	}
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func GetDb(dbName string) *sql.DB {
	db, ok := mapDB.Load(dbName)
	if !ok {
		panic(errors.New("no db exists"))
	}
	return db.(*sql.DB)
}
