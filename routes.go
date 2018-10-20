package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	var b passwordEntry
	e := c.Bind(&b)

	if e != nil {
		log.Fatal(e)
	}

	if ok && matches {
		existingPwEntries := userDataDB[uname]
		userDataDB[uname] = append(existingPwEntries, b)
	}
}

func GetData(c *gin.Context) {
	// Authenticate the user and then give their personal data
	uname, _, ok, matches := auth(c)

	if ok && matches {
		userEntries := userDataDB[uname]
		c.JSON(http.StatusOK, userEntries)
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
