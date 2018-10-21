package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

//
type passwordEntry struct {
	EntryUsername string `json:"entry_username"`
	EntryPassword string `json:"entry_password"`
	EntryDomain   string `json:"entry_domain"`
}

type authType struct {
	Email string `json:"email"`
	Pw    string `json:"pw"`
}

type userDataKey struct {
	KeyUsername string `json:"key_username"`
	KeyDomain   string `json:"key_domain"`
}

var (
	userExistsDB = make(map[string]string)
	userDataDB   = make(map[string]map[userDataKey]passwordEntry)
	userMetaData = make(map[string]string) // Actual name map
)

func main() {
	r := gin.Default()
	r.POST("/getdata", GetData)
	r.POST("/createuser", CreateUser)
	r.POST("/add", AddEntry)
	r.POST("/entry/delete", DeleteEntry)
	r.POST("/auth", Authenticate)
	err := r.RunTLS(":443", "ssl/dev.gganley.com.crt", "ssl/dev.gganley.com.key")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
