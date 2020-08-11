package main

import (
	"fmt"
	"github.com/yn295636/MyGoPractice/db"
	pb "github.com/yn295636/MyGoPractice/proto/greeter_service"
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
		"username":     in.Username,
		"gender":       in.Gender,
		"age":          in.Age,
		"external_uid": in.ExternalUid,
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

func QueryUserByUid(in *pb.GetUserFromDbRequest) (*pb.GetUserFromDbReply, error) {
	dbConn := db.GetDb(db.DbMyGoPractice)
	sqlStat := fmt.Sprintf("SELECT username,gender,age,external_uid FROM %v WHERE id=?", db.TblUser)
	row := dbConn.QueryRow(sqlStat, in.Uid)
	var (
		username    string
		gender      int32
		age         int32
		externalUid int32
	)
	err := row.Scan(&username, &gender, &age, &externalUid)
	if err != nil {
		log.Printf("Query from %v failed, %v", db.TblUser, err)
		return nil, err
	}

	var out = &pb.GetUserFromDbReply{
		Username:    username,
		Gender:      gender,
		Age:         age,
		ExternalUid: externalUid,
	}
	return out, nil
}
