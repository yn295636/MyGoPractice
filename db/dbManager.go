package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type TDbVersion int

const (
	DbMyGoPractice            = "mygopractice"
	TblMeta                   = "meta"
	TblUser                   = "user"
	TblPhone                  = "phone"
	MetaSchema                = "version INT(32) NOT NULL"
	DbVer1         TDbVersion = 1
	DbVer2         TDbVersion = 2
)

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

func createTableInTx(dbTx *sql.Tx, tableName, tableSchema string) error {
	sqlStat := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v(\n%v\n)\nCOLLATE='utf8mb4_unicode_ci'\nENGINE=InnoDB",
		tableName, tableSchema)
	_, err := dbTx.Exec(sqlStat)
	return err
}

func alterTableInTx(dbTx *sql.Tx, tableName, modifiedSchema string) error {
	sqlStat := fmt.Sprintf("ALTER TABLE %v \n%v",
		tableName, modifiedSchema)
	_, err := dbTx.Exec(sqlStat)
	return err
}

func getCurrentDbVersion(dbConn *sql.DB) (TDbVersion, error) {
	var count uint
	sqlStat1 := fmt.Sprintf("SELECT COUNT(version) FROM %v", TblMeta)
	row1 := dbConn.QueryRow(sqlStat1)
	err := row1.Scan(&count)
	if count == 0 {
		sqlStat2 := fmt.Sprintf("INSERT INTO %v (version) VALUES (?)", TblMeta)
		_, err = dbConn.Exec(sqlStat2, 0)
		if err != nil {
			log.Printf("Set db version to 0 failed, %v", err)
			return 1000000, err
		}
		return 0, nil
	}

	var result TDbVersion
	sqlState := fmt.Sprintf("SELECT version FROM %v", TblMeta)
	row := dbConn.QueryRow(sqlState)
	err = row.Scan(&result)
	if err != nil {
		log.Printf("Query db version failed, %v", err)
		return 1000000, err
	}
	return result, nil
}

func updateDbVersion(dbTx *sql.Tx, ver TDbVersion) error {
	sqlStat := fmt.Sprintf("UPDATE %v SET version=?", TblMeta)
	_, err := dbTx.Exec(sqlStat, ver)
	return err
}

func updateMyGoPracticeToVer1(dbConn *sql.DB) error {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Printf("DbMyGoPractice Tx begin error, %v", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	userSchema := `id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			username VARCHAR(16) NOT NULL DEFAULT '',
			PRIMARY KEY (id)`
	err = createTableInTx(tx, TblUser, userSchema)
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
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	phoneSchema := `id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			phone VARCHAR(16) NOT NULL DEFAULT '',
			uid BIGINT(20) UNSIGNED NOT NULL,
			PRIMARY KEY (id)`
	err = createTableInTx(tx, TblPhone, phoneSchema)
	if err != nil {
		log.Printf("Create table phone failed, %v", err)
		return err
	}

	userAddSchema := `ADD COLUMN gender TINYINT(3) NOT NULL DEFAULT 0 COMMENT '0:unknown, 1:male, 2:female',
			ADD COLUMN age TINYINT(3) NOT NULL DEFAULT 0`
	err = alterTableInTx(tx, TblUser, userAddSchema)
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
	defer dbMyGoPractice.Close()

	err = createTable(dbMyGoPractice, TblMeta, MetaSchema)
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

	return nil
}
