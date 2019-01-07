package qnsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	// 3rd
	_ "github.com/go-sql-driver/mysql"
)

// DatabaseCheck check database if not exist create table
func DatabaseCheck() bool {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	db, err := sql.Open("mysql", "root:@/aecm?charset=utf8")
	if nil != err {
		log.Fatal(err)
		return false
	}
	log.Printf("Database open success!")
	// create table mobile if not exists.
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `mobile`(`id` INT UNSIGNED AUTO_INCREMENT, `osVersion` VARCHAR(20) NOT NULL, `brand` VARCHAR(40) NOT NULL, `model` VARCHAR(40) NOT NULL, `sdkVersion` VARCHAR(20), `packageName` VARCHAR(100) NOT NULL, `author` VARCHAR(40), `instert_time` DATETIME DEFAULT NOW(), PRIMARY KEY ( `id` ))ENGINE=InnoDB DEFAULT CHARSET=utf8;")
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
func AddMobile(db *sql.DB, osVersion string, brand string, model string, sdkVersion string, packageName string, author string) bool {
	if db == nil {
		log.Println("sql.db is nil")
		return false
	}
	sqlStr := fmt.Sprintf("INSERT INTO mobile (osVersion, brand, model, sdkVersion, packageName, author) VALUES (\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\")",
		osVersion,
		brand,
		model,
		sdkVersion,
		packageName,
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
	osVersion   string
	brand       string
	model       string
	sdkVersion  string
	packageName string
	author      string
	insertTime  string
}

func struct2map(obj *aecmMobile) map[string]string {
	var data = make(map[string]string)
	data["osVersion"] = obj.osVersion
	data["brand"] = obj.brand
	data["model"] = obj.model
	data["sdkVersion"] = obj.sdkVersion
	data["packageName"] = obj.packageName
	data["author"] = obj.author
	data["insertTime"] = obj.insertTime
	return data
}

// QueryMobiles return json strings
func QueryMobiles(db *sql.DB, model string) string {
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
	if model != "" {
		sqlStr = fmt.Sprintf("SELECT * FROM aecm.mobile WHERE model='%s'", model)
	}
	query, err := db.Query(sqlStr)
	if err != nil {
		log.Println(err)
		return ""
	}
	mobileMap := make(map[int]map[string]string)
	for query.Next() {
		var id int
		var osVersion, brand, model, sdkVersion, packageName, author, insertTime string
		err := query.Scan(&id, &osVersion, &brand, &model, &sdkVersion, &packageName, &author, &insertTime)
		if err != nil {
			log.Fatal(err)
			continue
		}
		log.Printf("Query scan result: id:%d, osVersion:%s, brand:%s, model:%s, sdkVersion:%s, packageName:%s, author:%s, insertTime:%s\n",
			id,
			osVersion,
			brand,
			model,
			sdkVersion,
			packageName,
			author,
			insertTime)
		mobileMap[id] = struct2map(&aecmMobile{osVersion, brand, model, sdkVersion, packageName, author, insertTime})
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
func DeleteMobile(db *sql.DB, model string) bool {
	if db == nil {
		log.Fatalln("db is nil")
		return false
	}
	_, err := db.Exec("DELETE FROM mobile WHERE model=?", model)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
