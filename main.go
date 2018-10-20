package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type passwordEntry struct {
	username string
	password string
	domain   string
}

var (
	userExistsDB = make(map[string]string)

	userDataDB = make(map[string][]passwordEntry)
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

func CreateUser(c *gin.Context) {
	// inserts uname and pw into `userExistsDB`

	// {uname: "gganley", pw: "hellokitty"}
	// uname := req.PostForm.Get("uname")

	uname := c.PostForm("uname")
	pw := c.PostForm("pw")

	_, ok := userExistsDB[uname]

	if ok {
		_, e := c.Writer.Write([]byte("Username alread exists"))
		if e != nil {
			log.Fatal(e)
		}
	} else {
		// TODO: Actually implement this
		userExistsDB[uname] = pw
		userDataDB[uname] = []passwordEntry{}
	}
}

func AddEntry(c *gin.Context) {
	uname, _, ok, matches := auth(c)
	entryUsername := c.PostForm("entry_username")
	entryPassword := c.PostForm("entry_password")
	entryDomain := c.PostForm("entry_domain")
	if ok && matches {
		userDataDB[uname] = append(userDataDB[uname], passwordEntry{entryUsername, entryPassword, entryDomain})
	}
}

func GetData(c *gin.Context) {
	// Authenticate the user and then give their personal data
	uname, _, ok, matches := auth(c)

	if ok && matches {
		userEntries := userDataDB[uname]
		for _, v := range userEntries {

			_, err := c.Writer.Write([]byte(strings.Join([]string{v.domain, v.username, v.password}, ",")))

			if err != nil {
				log.Panic(err)
			}
		}
	}
}

// Returns the uname and pw, ok means that the entry exists and matches means that the pw is valid
func auth(c *gin.Context) (uname string, pw string, ok bool, matches bool) {
	uname = c.PostForm("uname")
	pw = c.PostForm("pw")
	val, ok := userExistsDB[uname]
	if val == pw {
		matches = true
	} else {
		matches = false
	}
	return
}
