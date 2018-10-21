package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func DeleteEntry(c *gin.Context) {

	var b struct {
		Auth  authType    `json:"auth"`
		Entry userDataKey `json:"entry"`
	}

	c.Bind(&b)

	ok, matches := authFunctino(b.Auth)

	if ok && matches {
		_, userDataExists := userDataDB[b.Auth.Email][b.Entry]

		if userDataExists {
			delete(userDataDB[b.Auth.Email], b.Entry)
			c.JSON(http.StatusOK, gin.H{"deleted": true})
		} else {
			c.JSON(http.StatusConflict, gin.H{"deleted": false})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Could not authenticate"})
	}
}

func CreateUser(c *gin.Context) {
	// inserts email and pw into `userExistsDB`

	// {email: "gganley", pw: "hellokitty"}
	// email := req.PostForm.Get("email")

	var b authType

	c.Bind(&b)

	fmt.Println("Creating user: ", b)
	ok, _ := authFunctino(b)

	if ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "User already exists"})
	} else {
		// TODO: Actually implement this
		userExistsDB[b.Email] = b.Pw
		userDataDB[b.Email] = make(map[userDataKey]passwordEntry)
		c.JSON(http.StatusOK, gin.H{"status": "User created"})
	}
}

func AddEntry(c *gin.Context) {

	type randomStruct struct {
		Auth  authType      `json:"auth"`
		Entry passwordEntry `json:"entry"`
	}

	var b randomStruct

	err := binding.JSON.Bind(c.Request, &b)
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Println(b)

	ok, matches := authFunctino(b.Auth)

	if ok && matches {
		userDataDB[b.Auth.Email][userDataKey{b.Entry.EntryUsername, b.Entry.EntryDomain}] = b.Entry
		c.JSON(http.StatusOK, gin.H{"added": true})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"status": "Unauthorized user"})
	}
}

func GetData(c *gin.Context) {
	// Authenticate the user and then give their personal data

	// Mild change
	var b authType

	binding.JSON.Bind(c.Request, &b)
	fmt.Println("Trying to get: ", b)

	ok, matches := authFunctino(b)

	if ok && matches {
		userEntries := userDataDB[b.Email]
		retVal := make(map[string]passwordEntry)
		for k, v := range userEntries {
			retVal[strings.Join([]string{k.KeyUsername, k.KeyDomain}, ":")] = v
		}
		c.JSON(http.StatusOK, retVal)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized user"})
	}
}

func Authenticate(c *gin.Context) {
	var b authType

	binding.JSON.Bind(c.Request, &b)

	ok, matches := authFunctino(b)

	if ok && matches {
		c.JSON(http.StatusOK, gin.H{"auth": true})
	} else {
		c.JSON(403, gin.H{"auth": false})
	}
}

// Returns the uname and pw, ok means that the entry exists and matches means that the pw is valid,
func authFunctino(auth authType) (ok bool, matches bool) {
	fmt.Println("trying to auth: ", auth, userExistsDB)
	val, ok := userExistsDB[auth.Email]

	if ok && val == auth.Pw {
		matches = true
	} else {
		fmt.Println(auth, val, ok)
		matches = false
	}
	return
}
