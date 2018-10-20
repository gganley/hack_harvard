package main

import (
	"log"
	// "fmt"
	// "io"
	"net/http"
	"strings"
)

type usernamePw struct {
	uname string
	pw    string
}

type passwordEntry struct {
	entryUsername string
	password      string
	domain        string
}

var (
	userExistsDB map[string]string

	userDataDB map[usernamePw][]passwordEntry
)

func main() {
	http.HandleFunc("/getdata", GetData)
	http.HandleFunc("/createuser", CreateUser)
	err := http.ListenAndServeTLS(":443", "ssl/dev.gganley.com.crt", "ssl/dev.gganley.com.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func CreateUser(writer http.ResponseWriter, req *http.Request) {
	// inserts uname and pw into `userExistsDB`

	// {uname: "gganley", pw: "hellokitty"}
	uname := req.PostForm.Get("uname")
	pw := req.PostForm.Get("pw")

	_, ok := userExistsDB[uname]

	if ok {
		_, e := writer.Write([]byte("Username alread exists"))
		if e != nil {
			log.Fatal(e)
		}
	} else {
		// TODO: Actually implement this
		userExistsDB[uname] = pw
		userDataDB[usernamePw{uname, pw}] = []passwordEntry{}
	}
}

func GetData(writer http.ResponseWriter, req *http.Request) {
	// Authenticate the user and then give their personal data
	uname := req.PostForm.Get("uname")
	pw := req.PostForm.Get("pw")

	userObject := usernamePw{uname, pw}
	val, ok := userExistsDB[uname]

	if ok && val == pw {
		userEntries := userDataDB[userObject]
		for _, v := range userEntries {

			_, err := writer.Write([]byte(strings.Join([]string{v.domain, v.entryUsername, v.password}, ",")))

			if err != nil {
				log.Panic(err)
			}
		}
	}
}
