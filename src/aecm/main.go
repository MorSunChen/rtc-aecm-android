package main

import (
	"aecm/http"
	"log"
)

func main() {
	// 指定 log 打印文件名 和 行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("Aecm Server starting...")
	http.StartHTTPServer("", 6688, "/v1/aecm/")
}
