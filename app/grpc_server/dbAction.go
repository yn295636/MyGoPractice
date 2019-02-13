package main

import (
	"github.com/yn295636/MyGoPractice/db"
	pb "github.com/yn295636/MyGoPractice/proto"
	"log"
)

func StoreUser(in *pb.StoreUserInDbRequest) (int64, error) {
	dbConn := db.GetDb(db.DbMyGoPractice)
	tx, err := dbConn.Begin()
	if err != nil {
		log.Printf("Begin tx failed, %v", err)
		return 0, err
	}
	defer db.TxPost(tx, err)

	kvMap := map[string]interface{}{
		"username": in.Username,
		"gender":   in.Gender,
		"age":      in.Age,
	}

	result, err := db.InsertIntoByTx(tx, db.TblUser, kvMap)
	if err != nil {
		log.Printf("Insert into %v failed, %v", db.TblUser, err)
		return 0, err
	}
	return result.LastInsertId()
}

func StorePhone(in *pb.StorePhoneInDbRequest) (int64, error) {
	dbConn := db.GetDb(db.DbMyGoPractice)
	tx, err := dbConn.Begin()
	if err != nil {
		log.Printf("Begin tx failed, %v", err)
		return 0, err
	}
	defer db.TxPost(tx, err)

	kvMap := map[string]interface{}{
		"phone": in.Phone,
		"uid":   in.Uid,
	}
	result, err := db.InsertIntoByTx(tx, db.TblPhone, kvMap)
	if err != nil {
		log.Printf("Insert into %v failed, %v", db.TblPhone, err)
		return 0, err
	}
	return result.LastInsertId()
}