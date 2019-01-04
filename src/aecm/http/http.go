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
		retStr := qnsql.QueryMobiles(sqlDb, "")
		if retStr != "" {
			fmt.Fprintf(w, retStr)
		} else {
			w.WriteHeader(501)
		}
	} else if r.RequestURI == "/v1/aecm/query" {
		model := r.Header.Get("model")
		retStr := qnsql.QueryMobiles(sqlDb, model)
		if retStr != "" {
			fmt.Fprintf(w, retStr)
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
func StartHTTPServer(ip string, port int, prifix string) {
	if !qnsql.DatabaseCheck() {
		log.Fatal("Database check failed, please check your mysql service!\n")
		return
	}
	sqlDb = qnsql.OpenDatabase()
	http.HandleFunc(prifix, handler)
	listenAddr := fmt.Sprintf("%s:%d", ip, port)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
