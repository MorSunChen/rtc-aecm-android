package qnsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	// 3rd
	_ "github.com/go-sql-driver/mysql"
)

func databaseCheck() bool {
	db, err := sql.Open("mysql", "root:Cgb815679@/aecm?charset=utf8")
	if nil != err {
		log.Fatal(err)
		return false
	}
	log.Printf("Database open success!")
	// create table mobile if not exists.
	_, err = db.Query("CREATE TABLE IF NOT EXISTS `mobile`(`id` INT UNSIGNED AUTO_INCREMENT, `mobile_id` VARCHAR(100) NOT NULL, `mobile_name` VARCHAR(40) NOT NULL, `author` VARCHAR(40) NOT NULL,`instert_time` DATETIME DEFAULT NOW(), PRIMARY KEY ( `id` ))ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	if err != nil {
		log.Println(err)
		return false
	}
	defer func(db *sql.DB) {
		if db != nil {
			db.Close()
			log.Println("Close sql.db.")
		}
	}(db)
	return true
}

// OpenDatabase open data base and return sql.DB pointer
func OpenDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:Cgb815679@/aecm?charset=utf8")
	if nil != err {
		log.Fatal(err)
		return nil
	}
	return db
}

// CloseDatabase close db
func CloseDatabase(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// AddMobile add new mobile record to sql
func AddMobile(db *sql.DB, mobileID string, mobileName string, author string) bool {
	if db == nil {
		log.Println("sql.db is nil")
		return false
	}
	sqlStr := fmt.Sprintf("INSERT INTO mobile (mobile_id, mobile_name, author) VALUES (\"%s\",\"%s\",\"%s\")",
		mobileID,
		mobileName,
		author)
	res, err := db.Exec(sqlStr)
	if nil != err {
		log.Fatal(err)
		return false
	}
	id, _ := res.LastInsertId()
	log.Println("last insert id:", id)
	return true
}

type aecmMobile struct {
	mobileID   string
	mobileName string
	author     string
	insertTime string
}

func struct2map(obj *aecmMobile) map[string]string {
	var data = make(map[string]string)
	data["mobileID"] = obj.mobileID
	data["mobileName"] = obj.mobileName
	data["author"] = obj.author
	data["insertTime"] = obj.insertTime
	return data
}

// QueryMobiles return json strings
func QueryMobiles(db *sql.DB, mobileID string) string {
	if db == nil {
		log.Println("sql.db is nil")
		return ""
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var sqlStr = "SELECT * FROM aecm.mobile"
	if mobileID != "" {
		sqlStr = fmt.Sprintf("SELECT * FROM aecm.mobile WHERE mobile_id='%s'", mobileID)
	}
	query, err := db.Query(sqlStr)
	if err != nil {
		log.Println(err)
		return ""
	}
	mobileMap := make(map[int]map[string]string)
	for query.Next() {
		var id int
		var mobileID, mobileName, author, insertTime string
		err := query.Scan(&id, &mobileID, &mobileName, &author, &insertTime)
		if err != nil {
			log.Fatal(err)
			continue
		}
		log.Printf("Query scan result: id:%d, mobileID:%s, mobileName:%s, author:%s, insertTime:%s\n",
			id,
			mobileID,
			mobileName,
			author,
			insertTime)
		mobileMap[id] = struct2map(&aecmMobile{mobileID, mobileName, author, insertTime})
	}
	jsonStr, err := json.Marshal(mobileMap)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	log.Printf("%s\n", jsonStr)
	return string(jsonStr)
}

// DeleteMobile delete this mobile from mysql
func DeleteMobile(db *sql.DB, mobile string) bool {
	if db == nil {
		log.Fatalln("db is nil")
		return false
	}
	_, err := db.Exec("DELETE FROM mobile WHERE mobile_id=?", mobile)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
