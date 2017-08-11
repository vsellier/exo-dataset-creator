package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
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

type SpaceActivity struct {
	SpaceId  int      `json:"id"`
	Activity Activity `json:"model"`
}

type Activity struct {
	Title string `json:"title"`
}

const (
	SESSION_URI          = "%s/rest/private/"
	USERS_URI            = "/rest/private/v1/social/users"
	SPACES_URI           = "/rest/private/v1/social/spaces"
	SPACE_ACTIVITIES_URI = "%s/rest/private/v1/social/spaces/%d/activities"
	USER_PREFIX          = "test"
	USER_PASSWORD        = "test123"
	USER_EMAIL           = "@test.com"
	SPACE_PREFIX         = "1000space"
	NB_USERS             = 1000
	NB_SPACES            = 2000
	NB_SPACES_ACTIVITIES = 10
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
		req.Header.Set("Content-Type", "application/json")
		fmt.Print("User ", name, " : ")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error : ", err)
		} else {
			fmt.Println(res.Status)
		}

		req.Body.Close()
		res.Body.Close()
	}
}

/*
curl -uroot:gtn -H'Content-Type: application/json' -X POST http://localhost:8080/rest/private/v1/social/spaces -d '{"displayName":"space1", "description":"space1", "visibility":"public", "subscription":"open"}'
*/
func createSpaces(h string, u string, p string) {
	for i := 0; i < NB_SPACES; i++ {
		name := fmt.Sprintf("%s%d", SPACE_PREFIX, i)
		newSpace := Space{DisplayName: name, Description: name, Visibility: "public", Subscription: "open"}
		json, _ := json.Marshal(newSpace)

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", h, SPACES_URI), bytes.NewBuffer(json))
		req.Header.Set("Content-Type", "application/json")
		fmt.Print("Space ", name, " : ")

		res, _ := client.Do(req)

		req.Body.Close()
		res.Body.Close()

		fmt.Println(res.Status)

	}
}

/*
 * Create NB_ACTIVITIES on each space
 */
func createSpacesActivities(h string, u string, p string) {

	for i := 1; i <= NB_SPACES_ACTIVITIES; i++ {
		for s := 1; s <= NB_SPACES; s++ {
			title := fmt.Sprintf("%s%d - Activity %d", SPACE_PREFIX, s, i)

			a := Activity{Title: title}

			fmt.Print(fmt.Sprintf("Creating activity %s ...", title))

			json, _ := json.Marshal(a)

			req, _ := http.NewRequest("POST", fmt.Sprintf(SPACE_ACTIVITIES_URI, h, s), bytes.NewBuffer(json))

			req.Close = true
			req.Header.Set("Content-Type", "application/json")

			res, _ := client.Do(req)
			fmt.Println(res.Status)

			req.Body.Close()
			res.Body.Close()

		}
	}
}

func getSession(h string, u string, p string) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(SESSION_URI, h), nil)

	req.Close = true
	req.SetBasicAuth(u, p)

	res, _ := client.Do(req)
	defer res.Body.Close()

	fmt.Println(client.Jar)

	fmt.Println("Session created " + res.Status)
}

var client http.Client
var sessionCookie string

func init() {

	jar, _ := cookiejar.New(&cookiejar.Options{})
	client = http.Client{Jar: jar}
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

	getSession(host, user, password)

	//	createUsers(host, user, password)
	//createSpaces(host, user, password)
	createSpacesActivities(host, user, password)

}
