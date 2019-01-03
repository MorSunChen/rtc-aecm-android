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

	if r.RequestURI == "/v1/aecm/add" {
		mobileID := r.Header.Get("mobileID")
		mobileName := r.Header.Get("mobileName")
		author := r.Header.Get("author")

		if qnsql.AddMobile(sqlDb, mobileID, mobileName, author) {
			fmt.Fprintf(w, "Sql operator result:success")
		} else {
			fmt.Fprintf(w, "Sql operator result:failed")
		}
	} else if r.RequestURI == "/v1/aecm/queryall" {
		retStr := qnsql.QueryMobiles(sqlDb, "")
		if retStr != "" {
			fmt.Fprintf(w, "Sql query success:%s", retStr)
		} else {
			fmt.Fprintf(w, "Sql query failed!")
		}
	} else if r.RequestURI == "/v1/aecm/query" {
		mobileID := r.Header.Get("mobileID")
		retStr := qnsql.QueryMobiles(sqlDb, mobileID)
		if retStr != "" {
			fmt.Fprintf(w, "Sql query success:%s", retStr)
		} else {
			fmt.Fprintf(w, "Sql query failed!")
		}
	} else if r.RequestURI == "/v1/aecm/delete" {
		mobileID := r.Header.Get("mobileID")
		if qnsql.DeleteMobile(sqlDb, mobileID) {
			fmt.Fprintf(w, "Sql delete success")
		} else {
			fmt.Fprintf(w, "Sql query failed!")
		}
	}
}

// StartHTTPServer Start http server
func StartHTTPServer(ip string, port int, prifix string) {
	sqlDb = qnsql.OpenDatabase()
	http.HandleFunc(prifix, handler)
	listenAddr := fmt.Sprintf("%s:%d", ip, port)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
