package qnsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	// 3rd
	_ "github.com/go-sql-driver/mysql"
)

type aecmMobile struct {
	osVersion   string
	brand       string
	model       string
	sdkVersion  string
	packageName string
	author      string
	insertTime  string
}

var mutex sync.RWMutex
var mobilesMap = make(map[int]aecmMobile) // key:ID; value:struct of aecmMobile

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

// DatabaseCheck check database if not exist create table
func DatabaseCheck(sqlUser string, sqlPwd string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	dataSourceName := fmt.Sprintf("%s:%s@/aecm?charset=utf8", sqlUser, sqlPwd)
	db, err := sql.Open("mysql", dataSourceName)
	if nil != err {
		log.Println(err)
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
func OpenDatabase(sqlUser string, sqlPwd string) *sql.DB {
	dataSourceName := fmt.Sprintf("%s:%s@/aecm?charset=utf8", sqlUser, sqlPwd)
	db, err := sql.Open("mysql", dataSourceName)
	if nil != err {
		log.Println(err)
		return nil
	}
	// init global varaibal mobilesMap
	updateMobilesMap(db)
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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	sqlStr := fmt.Sprintf("INSERT INTO mobile (osVersion, brand, model, sdkVersion, packageName, author) VALUES (\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\")",
		osVersion,
		brand,
		model,
		sdkVersion,
		packageName,
		author)
	res, err := db.Exec(sqlStr)
	if nil != err {
		log.Println(err)
		return false
	}
	id, _ := res.LastInsertId()
	log.Println("last insert id:", id)

	updateMobilesMap(db)

	return true
}

// updateMobilesMap return json strings
func updateMobilesMap(db *sql.DB) {
	if db == nil {
		log.Println("sql.db is nil")
		return
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var sqlStr = "SELECT * FROM aecm.mobile"
	query, err := db.Query(sqlStr)
	if err != nil {
		log.Println(err)
		return
	}
	mutex.Lock()
	for query.Next() {
		var id int
		var osVersion, brand, model, sdkVersion, packageName, author, insertTime string
		err := query.Scan(&id, &osVersion, &brand, &model, &sdkVersion, &packageName, &author, &insertTime)
		if err != nil {
			log.Println(err)
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
		mobilesMap[id] = aecmMobile{osVersion, brand, model, sdkVersion, packageName, author, insertTime}
	}
	mutex.Unlock()
}

// QueryMobiles query mobiles from memory, and return json string
func QueryMobiles(model string) (string, bool) {
	mutex.RLock()
	if model == "" {
		tmpMap := make(map[int]map[string]string)
		for key, val := range mobilesMap {
			tmpMap[key] = struct2map(&val)
		}
		jsonStr, err := json.Marshal(tmpMap)
		if err != nil {
			log.Println(err)
			mutex.RUnlock()
			return "", false
		}
		log.Printf("%s\n", jsonStr)
		mutex.RUnlock()
		return string(jsonStr), true
	} else {
		for _, val := range mobilesMap {
			if val.model == model {
				mutex.RUnlock()
				return "", true
			}
		}
		mutex.RUnlock()
		return "", false
	}
}

// DeleteMobile delete this mobile from mysql
func DeleteMobile(db *sql.DB, model string) bool {
	if db == nil {
		log.Println("db is nil")
		return false
	}
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	_, err := db.Exec("DELETE FROM mobile WHERE model=?", model)
	if err != nil {
		log.Println(err)
		return false
	}
	updateMobilesMap(db)
	return true
}
