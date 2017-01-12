package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type Space struct {
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	Visibility   string `json:"visibility"`
	Subscription string `json:"subscription"`
}

const (
	USERS_URI     = "/rest/private/v1/social/users"
	SPACES_URI    = "/rest/private/v1/social/spaces"
	USER_PREFIX   = "test"
	USER_PASSWORD = "test123"
	USER_EMAIL    = "@test.com"
	SPACE_PREFIX  = "space"
	NB_USERS      = 1000
	NB_SPACES     = 100
)

func usage(arguments []string) {
	fmt.Println(fmt.Sprintf("%s <base url> <user> <password>", arguments[0]))
	os.Exit(1)
}

func createUsers(h string, u string, p string) {
	for i := 0; i < NB_USERS; i++ {
		name := fmt.Sprintf("%s%d", USER_PREFIX, i)
		newUser := User{Username: name, Password: USER_PASSWORD, Email: fmt.Sprintf("%s%s", name, USER_EMAIL), Firstname: name, Lastname: name}
		json, _ := json.Marshal(newUser)
		// fmt.Println(string(json))

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", h, USERS_URI), bytes.NewBuffer(json))
		req.SetBasicAuth(u, p)
		req.Header.Set("Content-Type", "application/json")
		fmt.Print("User ", name, " : ")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error : ", err)
		} else {
			fmt.Println(res.Status)
		}
	}
}

func createSpaces(h string, u string, p string) {
	for i := 0; i < NB_SPACES; i++ {
		name := fmt.Sprintf("%s%d", SPACE_PREFIX, i)
		newSpace := Space{DisplayName: name, Description: name, Visibility: "public", Subscription: "open"}
		json, _ := json.Marshal(newSpace)

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", h, SPACES_URI), bytes.NewBuffer(json))
		req.SetBasicAuth(u, p)
		req.Header.Set("Content-Type", "application/json")
		fmt.Print("Space ", name, " : ")

		res, _ := client.Do(req)
		fmt.Println(res.Status)

	}
}

var client http.Client

func init() {
	client = http.Client{}
}

func main() {
	arguments := os.Args

	if len(arguments) != 4 {
		usage(arguments)
	}

	host := arguments[1]
	user := arguments[2]
	password := arguments[3]

	fmt.Println(fmt.Sprintf("Using host %s and user %s", host, user))

	createUsers(host, user, password)
	createSpaces(host, user, password)

}
