package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeleteEntry(c *gin.Context) {
	auth, ok, matches := auth(c)

	var b struct {
		authType    `json:"auth"`
		userDataKey `json:"entry"`
	}
	c.Bind(&b)
	if ok && matches {
		_, userDataExists := userDataDB[auth.Email][b.userDataKey]

		if userDataExists {
			delete(userDataDB[auth.Email], b.userDataKey)
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

	auth, ok, _ := auth(c)

	if ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "User already exists"})
	} else {
		// TODO: Actually implement this
		userExistsDB[auth.Email] = auth.Pw
		userDataDB[auth.Email] = make(map[userDataKey]passwordEntry)
		c.JSON(http.StatusOK, gin.H{"status": "User created"})
	}
}

func AddEntry(c *gin.Context) {
	auth, ok, matches := auth(c)
	var b struct {
		Auth  authType      `json:"auth"`
		Entry passwordEntry `json:"entry"`
	}

	c.Bind(&b)

	if ok && matches {
		userDataDB[auth.Email][userDataKey{b.Entry.EntryUsername, b.Entry.EntryDomain}] = b.Entry
		c.JSON(http.StatusOK, gin.H{"added": true})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized user"})
	}
}

func GetData(c *gin.Context) {
	// Authenticate the user and then give their personal data
	auth, ok, matches := auth(c)

	if ok && matches {
		userEntries := userDataDB[auth.Email]
		retVal := make(map[string]passwordEntry)
		for k, v := range userEntries {
			retVal[strings.Join([]string{k.KeyUsername, k.KeyDomain}, ":")] = v
		}
		c.JSON(http.StatusOK, retVal)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized user"})
	}
}

// Returns the uname and pw, ok means that the entry exists and matches means that the pw is valid
func auth(c *gin.Context) (auth authType, ok bool, matches bool) {
	c.Bind(&auth)
	val, ok := userExistsDB[auth.Email]
	if ok && val == auth.Pw {
		matches = true
	} else {
		matches = false
	}
	return
}
