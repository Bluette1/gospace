// Copyright 2011 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"crypto/rand"
	"net/http"
	"os"
	"io"
	"strings"
	"log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/url"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"
)
const shortDuration = 10000 * time.Millisecond //revert to 1ms

var (
	service  *oauth1a.Service
	sessions map[string]*oauth1a.UserConfig
)
//Struct to parse webhook load
type WebhookLoad struct {
	UserId           string  `json:"for_user_id"`
	TweetCreateEvent []Tweet `json:"tweet_create_events"`
}

//Struct to parse tweet
type Tweet struct {
	Id    int64
	IdStr string `json:"id_str"`
	User  User
	Text  string
}

//Struct to parse user
type User struct {
	Id     int64
	IdStr  string `json:"id_str"`
	Name   string
	Handle string `json:"screen_name"`
}

func NewSessionID() string {
	c := 128
	b := make([]byte, c)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		panic("Could not generate random number")
	}
	return base64.URLEncoding.EncodeToString(b)
}

func GetSessionID(req *http.Request) (id string, err error) {
	var c *http.Cookie
	if c, err = req.Cookie("session_id"); err != nil {
		return
	}
	id = c.Value
	return
}

func SessionStartCookie(id string) *http.Cookie {
	return &http.Cookie{
		Name:   "session_id",
		Value:  id,
		MaxAge: 60,
		Secure: false,
		Path:   "/",
	}
}

func SessionEndCookie() *http.Cookie {
	return &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: 0,
		Secure: false,
		Path:   "/",
	}
}

func main(){
	// Load env
	sessions = map[string]*oauth1a.UserConfig{}

	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file")
			fmt.Println("Error loading .env file")
	}
	fmt.Println("Starting Server")

	service = &oauth1a.Service{
		RequestURL:   "https://api.twitter.com/oauth/request_token",
		AuthorizeURL: "https://api.twitter.com/oauth/authorize",
		AccessURL:    "https://api.twitter.com/oauth/access_token",
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    os.Getenv("CONSUMER_KEY"),
			ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
			CallbackURL:    os.Getenv("APP_URL") + "/callback/",
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}
	
	//Check for -register in agument list
	if args := os.Args; len(args) > 1 && args[1] == "-register"{
			go registerWebhook()
	}

	if args := os.Args; len(args) > 1 && args[1] == "-subscribe"{
		go subscribeWebhook()
}

if args := os.Args; len(args) > 1 && args[1] == "-delete"{
	go deleteWebhook()
}

if args := os.Args; len(args) > 1 && args[1] == "-reply"{
	go ReplyToTweet(args[2], args[3])
}
if args := os.Args; len(args) > 1 && args[1] == "-send"{
  twtStr := strings.Join(args[2:], " ")
	go SendTweet(twtStr)
}

	//Create a netw Mux Handler
	m := mux.NewRouter()
	//Listen to the base url and send a response
// 	m.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
// 		writer.WriteHeader(200)
// 		fmt.Fprintf(writer, "Server is up and running")
// })
	m.HandleFunc("/", BaseHandler)
	//Listen to crc check and handle
	m.HandleFunc("/webhook/twitter", CrcCheck).Methods("GET")
	 //Listen to webhook event and handle
	 m.HandleFunc("/webhook/twitter", WebhookHandler).Methods("POST")

	m.HandleFunc("/signin/", SignInHandler)
	m.HandleFunc("/callback/", CallbackHandler)

	//Start Server
	server := &http.Server{
			Handler: m,
	}
	server.Addr = ":9090"
	server.ListenAndServe()
}

func ReplyToTweet(tweet string, reply_id string) (*Tweet, error) {
	fmt.Println("Sending tweet as reply to " + reply_id)
	//Initialize tweet object to store response in
	var responseTweet Tweet
	//Add params
	params := url.Values{}
	//These may be irrelevant because we use query params
	params.Set("status", tweet)
	params.Set("in_reply_to_status_id",reply_id)
	//Grab client and post
	var (
		err    error
		client *twittergo.Client
		req    *http.Request
		resp   *twittergo.APIResponse
	)
	client, _ = LoadCredentials()
	if err != nil {
		fmt.Printf("Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}
	path := "/1.1/statuses/update.json?status=" + tweet + "&in_reply_to_status_id=" + reply_id
	
	body := strings.NewReader(params.Encode())
	req, err = http.NewRequest("POST", path, body)
	if err != nil {
			return nil, err
	}
	defer req.Body.Close()
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	//Decode response and send out
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	err = json.Unmarshal(respBody, &responseTweet)
	if err != nil{
			return  nil,err
	}
	return &responseTweet, nil
}

func SendTweet(tweet string) (*Tweet, error)  {
fmt.Println("Sending tweet... " )
	// Initialize tweet object to store response in
	var responseTweet Tweet
	//Add params
	params := url.Values{}
	//These may be irrelevant because we use query params
	params.Set("status", tweet)
	//Grab client and post
	var (
		err    error
		client *twittergo.Client
		req    *http.Request
		resp   *twittergo.APIResponse
	)
	client, _ = LoadCredentials()
	if err != nil {
		fmt.Printf("Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}
	path := "/1.1/statuses/update.json?status=" + tweet
	
	body := strings.NewReader(params.Encode())
	req, err = http.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	//Decode response and send out
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	err = json.Unmarshal(respBody, &responseTweet)
	if err != nil{
		return  nil,err
	}
	return &responseTweet, nil
}

func WebhookHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Handler called")
    //Read the body of the tweet
    body, _ := ioutil.ReadAll(request.Body)
    //Initialize a webhok load obhject for json decoding
    var load WebhookLoad
    err := json.Unmarshal(body, &load)
    if err != nil {
			fmt.Println("An error occured: " + err.Error())
    }
    //Check if it was a tweet_create_event and tweet was in the payload and it was not tweeted by the bot
    //Should NOT send reply tweets to oneself // as it creates a chain of
		//duplicate tweets
		if len(load.TweetCreateEvent) < 1 || load.UserId == load.TweetCreateEvent[0].User.IdStr { 
      return
    }
    //Send `So true...` as a reply to the tweet, replies need to begin with the handles
    _, err = ReplyToTweet("@"+load.TweetCreateEvent[0].User.Handle+" So true...", load.TweetCreateEvent[0].IdStr)
    if err != nil {
			fmt.Println("An error occured:")
			fmt.Println(err.Error())
    } else{
      fmt.Println("Tweet sent successfully")
    }
}

func CrcCheck(writer http.ResponseWriter, request *http.Request){
	//Set response header to json type
	writer.Header().Set("Content-Type", "application/json")
	//Get crc token in parameter
	token := request.URL.Query()["crc_token"]
	if len(token) < 1 {
		fmt.Fprintf(writer,"No crc_token given")
		return
	}

	//Encrypt and encode in base 64 then return
	h := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	h.Write([]byte(token[0]))
	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
	//Generate response string map
	response := make(map[string]string)
	response["response_token"] =  "sha256=" + encoded
	//Turn response map to json and send it to the writer
	responseJson, _ := json.Marshal(response)
	// fmt.Println(string(responseJson))
	fmt.Fprintf(writer, string(responseJson))
}

func LoadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	if err != nil {
		return
	}
	lines := strings.Split(string(credentials), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    lines[0],
		ConsumerSecret: lines[1],
	}
	user := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	client = twittergo.NewClient(config, user)
	return
}
func registerWebhook(){
	fmt.Println("Registering webhook...")
	var (
		err    error
		client *twittergo.Client
		req    *http.Request
		resp   *twittergo.APIResponse
	)
	client, err = LoadCredentials()
	if err != nil {
		fmt.Printf("Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}

	//Set parameters
	values := url.Values{}
	appUrl := url.QueryEscape(os.Getenv("APP_URL")+"/webhook/twitter")
	path := "/1.1/account_activity/all/"+os.Getenv("WEBHOOK_ENV")+"/webhooks.json?url="+appUrl

	values.Set("url", appUrl)

	body := strings.NewReader(values.Encode())
	req, err = http.NewRequest("POST", path, body)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	defer req.Body.Close()
	
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}

	// Parse response and check response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Could not parse response: %v\n", err)
		os.Exit(1)
	}
  fmt.Println(string(respBody))
}

func subscribeWebhook(){
	fmt.Println("Subscribing webapp...")
	var (
		err    error
		client *twittergo.Client
		req    *http.Request
		resp   *twittergo.APIResponse
	)
	client, _ = LoadCredentials()
	if err != nil {
		fmt.Printf("Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}

	path := "/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/subscriptions.json"
	values := url.Values{}
	payload := strings.NewReader(values.Encode())
	req, _ = http.NewRequest("POST", path, payload)
	defer req.Body.Close()

	//If response code is 204 it was successful
	resp, _ = client.SendRequest(req)
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 204 {
		fmt.Println("Subscribed successfully")
	} else if resp.StatusCode!= 204 {
		fmt.Println("Could not subscribe the webhook. Response below:")
		fmt.Println(string(body))
	}
}

func deleteWebhook(){
	fmt.Println("Deleting webhook...")
	var (
		err    error
		client *twittergo.Client
		req    *http.Request
		resp   *twittergo.APIResponse
	)
	client, err = LoadCredentials()
	if err != nil {
		fmt.Printf("Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}

	//Set parameters
	values := url.Values{}
  path := "/1.1/account_activity/all/"+os.Getenv("WEBHOOK_ENV")+"/webhooks/" + os.Getenv("WEBHOOK_ID") + ".json"

	body := strings.NewReader(values.Encode())
	req, err = http.NewRequest("DELETE", path, body)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	defer req.Body.Close()
	
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}

	// Parse response and check response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Could not parse response: %v\n", err)
		os.Exit(1)
	}
	if resp.StatusCode == 204 {
		fmt.Println("Webhook deleted successfully")
  } else if resp.StatusCode!= 204 {
		fmt.Println("Could not delete the webhook. Response below:")
		fmt.Println(string(respBody))
  }
}


func SignInHandler(rw http.ResponseWriter, req *http.Request) {
	var (
		url       string
		err       error
		sessionID string
	)
	httpClient := new(http.Client)
	userConfig := &oauth1a.UserConfig{}
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
  defer cancel()
	if err = userConfig.GetRequestToken(ctx, service, httpClient); err != nil {
		log.Printf("Could not get request token: %v", err)
		http.Error(rw, "Problem getting the request token", 500)
		return
	}
	if url, err = userConfig.GetAuthorizeURL(service); err != nil {
		log.Printf("Could not get authorization URL: %v", err)
		http.Error(rw, "Problem getting the authorization URL", 500)
		return
	}
	log.Printf("Redirecting user to %v\n", url)
	sessionID = NewSessionID()
	log.Printf("Starting session %v\n", sessionID)
	sessions[sessionID] = userConfig
	http.SetCookie(rw, SessionStartCookie(sessionID))
	http.Redirect(rw, req, url, http.StatusFound)
}

func CallbackHandler(rw http.ResponseWriter, req *http.Request) {
	var (
		err        error
		token      string
		verifier   string
		sessionID  string
		userConfig *oauth1a.UserConfig
		ok         bool
	)
	log.Printf("Callback hit. %v current sessions.\n", len(sessions))

	if sessionID, err = GetSessionID(req); err != nil {
		log.Printf("Got a callback with no session id: %v\n", err)
		http.Error(rw, "No session found", 400)
		return
	}
	if userConfig, ok = sessions[sessionID]; !ok {
		log.Printf("Could not find user config in sesions storage.")
		http.Error(rw, "Invalid session", 400)
		return
	}
	if token, verifier, err = userConfig.ParseAuthorize(req, service); err != nil {
		log.Printf("Could not parse authorization: %v", err)
		http.Error(rw, "Problem parsing authorization", 500)
		return
	}
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	httpClient := new(http.Client)
	if err = userConfig.GetAccessToken(ctx, token, verifier, service, httpClient); err != nil {
		log.Printf("Error getting access token: %v", err)
		http.Error(rw, "Problem getting an access token", 500)
		return
	}
	log.Printf("Ending session %v.\n", sessionID)
	delete(sessions, sessionID)
	http.SetCookie(rw, SessionEndCookie())
	rw.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprintf(rw, "<pre>")
	fmt.Fprintf(rw, "Access Token: %v\n", userConfig.AccessTokenKey)
	fmt.Fprintf(rw, "Token Secret: %v\n", userConfig.AccessTokenSecret)
	fmt.Fprintf(rw, "Screen Name:  %v\n", userConfig.AccessValues.Get("screen_name"))
	fmt.Fprintf(rw, "User ID:      %v\n", userConfig.AccessValues.Get("user_id"))
	fmt.Fprintf(rw, "</pre>")
	fmt.Fprintf(rw, "<a href=\"/signin\">Sign in again</a>")
}

func BaseHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprintf(rw, "<a href=\"/signin\">Sign in</a>")
}