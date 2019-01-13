package http

import (
	"aecm/qnsql"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var sqlDb *sql.DB

func handler(w http.ResponseWriter, r *http.Request) {

	log.Println("New request addr and uri is:", r.RemoteAddr, r.RequestURI)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if r.RequestURI == "/v1/aecm/add" {
		osVersion := r.Header.Get("osVersion")
		brand := r.Header.Get("brand")
		model := r.Header.Get("model")
		sdkVersion := r.Header.Get("sdkVersion")
		packageName := r.Header.Get("packageName")
		author := r.Header.Get("author")

		if !qnsql.AddMobile(sqlDb, osVersion, brand, model, sdkVersion, packageName, author) {
			w.WriteHeader(501)
		}
	} else if r.RequestURI == "/v1/aecm/queryall" {
		retStr, _ := qnsql.QueryMobiles("")
		if retStr != "" {
			w.Header().Set("Content-Type", "text/json; charset=utf-8")
			fmt.Fprintf(w, retStr)
		} else {
			w.WriteHeader(501)
		}
	} else if r.RequestURI == "/v1/aecm/query" {
		model := r.Header.Get("model")
		_, ret := qnsql.QueryMobiles(model)
		if ret == true {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(501)
		}
	} else if r.RequestURI == "/v1/aecm/delete" {
		model := r.Header.Get("model")
		if !qnsql.DeleteMobile(sqlDb, model) {
			w.WriteHeader(501)
		}
	}
}

// StartHTTPServer Start http server
func StartHTTPServer(addr string, prifix string, sqlUser string, sqlPwd string) {
	if !qnsql.DatabaseCheck(sqlUser, sqlPwd) {
		log.Fatal("Database check failed, please check your mysql service!\n")
		return
	}
	sqlDb = qnsql.OpenDatabase(sqlUser, sqlPwd)
	http.HandleFunc(prifix, handler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
