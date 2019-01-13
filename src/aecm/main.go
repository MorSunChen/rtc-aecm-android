package main

import (
	"aecm/http"
	"flag"
	"log"
)

func main() {
	addrStr := flag.String("addr", "0.0.0.0:6688", "ip address and port, example: 0.0.0.0:6688")
	prefixStr := flag.String("prefix", "/v1/aecm/", "api prefix string")
	sqlUser := flag.String("user", "root", "mysql user name")
	sqlPwd := flag.String("password", "", "mysql password")
	flag.Parse()

	log.Printf("Inpout params: addr:%s, prefix:%s, sql user:%s, sql pwd:%s\n",
		*addrStr,
		*prefixStr,
		*sqlUser,
		*sqlPwd)

	// 指定 log 打印文件名 和 行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("Aecm Server starting...")

	// Start server
	http.StartHTTPServer(*addrStr, *prefixStr, *sqlUser, *sqlPwd)
}
