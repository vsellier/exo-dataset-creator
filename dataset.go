package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"math/rand"
	"sync"
	"time"
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

type SpaceMembership struct {
	User  string `json:"user"`
	Space string `json:"space"`
}

const (
	SESSION_URI           = "%s/rest/private/"
	USERS_URI             = "/rest/private/v1/social/users"
	SPACES_URI            = "/rest/private/v1/social/spaces"
	SPACE_ACTIVITIES_URI  = "%s/rest/private/v1/social/spaces/%d/activities"
	SPACE_MEMBERSHIP_URL  = "%s/rest/private/v1/social/spacesMemberships"
	USER_PREFIX           = "test"
	USER_PASSWORD         = "test123"
	USER_EMAIL            = "@test.com"
	SPACE_PREFIX          = "space"
	NB_USERS              = 100
	NB_SPACES             = 2000
	NB_SPACES_ACTIVITIES  = 10
	SPACE_ACTIVITY_LENGTH = 200
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 !@#$%^&*()_+-=[]\\{}|;':\",./?"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func usage(arguments []string) {
	fmt.Println(fmt.Sprintf("%s <base url> <user> <password>", arguments[0]))
	os.Exit(1)
}

func createUser(wg *sync.WaitGroup, c <-chan string, h string) {
	for {
		name := <-c

		t0 := time.Now()

		newUser := User{Username: name, Password: USER_PASSWORD, Email: fmt.Sprintf("%s%s", name, USER_EMAIL), Firstname: name, Lastname: name}
		json, _ := json.Marshal(newUser)

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", h, USERS_URI), bytes.NewBuffer(json))
		req.Header.Set("Content-Type", "application/json")

		res, _ := client.Do(req)

		req.Body.Close()
		res.Body.Close()

		t1 := time.Now()

		fmt.Println("User ", name, " : "+res.Status, " in ", t1.Sub(t0))

		wg.Done()
	}
}

func createUsers(h string) {
	var wg sync.WaitGroup

	c := make(chan string)

	////-------------
	// Create any go routines here
	go createUser(&wg, c, h)
	go createUser(&wg, c, h)
	go createUser(&wg, c, h)
	go createUser(&wg, c, h)
	go createUser(&wg, c, h)
	go createUser(&wg, c, h)

	t0 := time.Now()

	for i := 0; i < NB_USERS; i++ {
		wg.Add(1)

		name := fmt.Sprintf("%s%d", USER_PREFIX, i)
		c <- name
	}
	fmt.Println("Waiting for the threads to finish")
	wg.Wait()
	t1 := time.Now()
	fmt.Println("All thread done")
	fmt.Println("Users created in ", t1.Sub(t0))

}

func addUserToSpace(h string, space string, user string) {
	t0 := time.Now()
	membership := SpaceMembership{User: user, Space: space}
	json, _ := json.Marshal(membership)

	req, _ := http.NewRequest("POST", fmt.Sprintf(SPACE_MEMBERSHIP_URL, h), bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	res, _ := client.Do(req)

	req.Body.Close()
	res.Body.Close()

	t1 := time.Now()

	fmt.Println("Add user", user, "to space", space, "in", t1.Sub(t0), res.Status)
}

func addUsersToSpaces(h string) {
	function := func(wg *sync.WaitGroup, c <-chan *SpaceMembership) {
		for {
			m := <-c

			addUserToSpace(h, m.Space, m.User)

			wg.Done()
		}
	}

	t0 := time.Now()

	var wg sync.WaitGroup
	// Updating spaces is not thread safe
	var channels []chan *SpaceMembership
	channels = append(channels, make(chan *SpaceMembership))
	channels = append(channels, make(chan *SpaceMembership))
	channels = append(channels, make(chan *SpaceMembership))
	channels = append(channels, make(chan *SpaceMembership))
	channels = append(channels, make(chan *SpaceMembership))

	for _, c := range channels {
		go function(&wg, c)
	}

	for sid := 0; sid < NB_SPACES; sid++ {
		for u := 0; u < NB_USERS; u++ {
			u := fmt.Sprintf("%s%d", USER_PREFIX, u)
			s := fmt.Sprintf("%s%d", SPACE_PREFIX, sid)
			m := SpaceMembership{User: u, Space: s}

			wg.Add(1)

			pos := sid % len(channels)
			channels[pos] <- &m
		}
	}

	fmt.Println("Waiting for the threads to finish")
	wg.Wait()
	t1 := time.Now()
	fmt.Println("All thread done")
	fmt.Println("User attacted to spaces in ", t1.Sub(t0))

}

func createSpace(wg *sync.WaitGroup, c <-chan int, h string) {
	// TODO cleary exit
	for {
		t0 := time.Now()
		id := <-c

		name := fmt.Sprintf("%s%d", SPACE_PREFIX, id)
		newSpace := Space{DisplayName: name, Description: name, Visibility: "public", Subscription: "open"}
		json, _ := json.Marshal(newSpace)

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", h, SPACES_URI), bytes.NewBuffer(json))
		req.Header.Set("Content-Type", "application/json")

		res, _ := client.Do(req)

		req.Body.Close()
		res.Body.Close()

		t1 := time.Now()

		fmt.Println("Space ", name, " : "+res.Status, " in ", t1.Sub(t0))

		wg.Done()
	}
}

/*
curl -uroot:gtn -H'Content-Type: application/json' -X POST http://localhost:8080/rest/private/v1/social/spaces -d '{"displayName":"space1", "description":"space1", "visibility":"public", "subscription":"open"}'
*/
func createSpaces(h string, u string, p string) {

	var wg sync.WaitGroup

	sc := make(chan int)

	////-------------
	// Create any go routines here
	// For the moment it's not possible
	// to create more than one space at a time
	// due to https://jira.exoplatform.org/browse/SOC-5697
	go createSpace(&wg, sc, h)

	t0 := time.Now()

	for i := 0; i < NB_SPACES; i++ {
		wg.Add(1)
		sc <- i
	}
	fmt.Println("Waiting for the threads to finish")
	wg.Wait()
	t1 := time.Now()
	fmt.Println("All thread done")
	fmt.Println("Spaces created in ", t1.Sub(t0))

}

/*
 * Create NB_ACTIVITIES on each space
 */
func createSpacesActivities(h string, u string, p string) {
	t0 := time.Now()

	for i := 1; i <= NB_SPACES_ACTIVITIES; i++ {
		for s := 1; s <= NB_SPACES; s++ {
			ta0 := time.Now()
			title := RandStringBytes(SPACE_ACTIVITY_LENGTH)

			a := Activity{Title: title}

			fmt.Print(fmt.Sprintf("Creating activity spaceId=%d actitivyCount=%d ...", s, i))

			json, _ := json.Marshal(a)

			req, _ := http.NewRequest("POST", fmt.Sprintf(SPACE_ACTIVITIES_URI, h, s), bytes.NewBuffer(json))

			req.Close = true
			req.Header.Set("Content-Type", "application/json")

			res, _ := client.Do(req)
			fmt.Print(res.Status)

			req.Body.Close()
			res.Body.Close()

			ta1 := time.Now()
			fmt.Println(" in ", ta1.Sub(ta0))

		}
	}
	t1 := time.Now()
	fmt.Println("Activities created in ", t1.Sub(t0))
}

func createSpacesActivitiy(h string, id int, content string) {
	ta0 := time.Now()

	a := Activity{Title: content}

	fmt.Print(fmt.Sprintf("Creating activity spaceId=%d ...", id))

	json, _ := json.Marshal(a)

	req, _ := http.NewRequest("POST", fmt.Sprintf(SPACE_ACTIVITIES_URI, h, id), bytes.NewBuffer(json))

	req.Close = true
	req.Header.Set("Content-Type", "application/json")

	res, _ := client.Do(req)
	fmt.Print(res.Status)

	req.Body.Close()
	res.Body.Close()

	ta1 := time.Now()
	fmt.Println(" in ", ta1.Sub(ta0))
}

func getSession(h string, u string, p string) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(SESSION_URI, h), nil)

	req.Close = true
	req.SetBasicAuth(u, p)

	res, _ := client.Do(req)
	defer res.Body.Close()

	//fmt.Println(client.Jar)

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

	createUsers(host)
	//createSpaces(host, user, password)
	addUsersToSpaces(host)
	//createSpacesActivities(host, user, password)
	//createSpacesActivitiy(host, 1, "<script language=\"javascript\">alert(\"test\")</script>test")

}
