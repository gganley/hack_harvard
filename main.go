package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type passwordEntry struct {
	EntryUsername string `json:"entry_username"`
	EntryPassword string `json:"entry_password"`
	EntryDomain   string `json:"entry_domain"`
}

var (
	userExistsDB = make(map[string]string)
	userDataDB   = make(map[string][]passwordEntry)
)

func main() {
	r := gin.Default()
	r.POST("/getdata", GetData)
	r.POST("/createuser", CreateUser)
	r.POST("/add", AddEntry)
	err := r.RunTLS(":443", "ssl/dev.gganley.com.crt", "ssl/dev.gganley.com.key")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
