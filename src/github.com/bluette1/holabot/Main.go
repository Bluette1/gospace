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
	"net/http"
	"os"
	"strings"
	"log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/url"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

func main(){
	// Load env
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file")
			fmt.Println("Error loading .env file")
	}
	fmt.Println("Starting Server")
	
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

	//Create a new Mux Handler
	m := mux.NewRouter()
	//Listen to the base url and send a response
	m.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(200)
			fmt.Fprintf(writer, "Server is up and running")
	})
	//Listen to crc check and handle
	m.HandleFunc("/webhook/twitter", CrcCheck).Methods("GET")
	 //Listen to webhook event and handle
	 m.HandleFunc("/webhook/twitter", WebhookHandler).Methods("POST")

	//Start Server
	server := &http.Server{
			Handler: m,
	}
	server.Addr = ":9090"
	server.ListenAndServe()
}
func SendTweet(tweet string, reply_id string) (*Tweet, error) {
	fmt.Println("Sending tweet as reply to " + reply_id)
	//Initialize tweet object to store response in
	var responseTweet Tweet
	//Add params
	params := url.Values{}
	params.Set("status",tweet)
	params.Set("in_reply_to_status_id",reply_id)
	//Grab client and post
	// client := CreateClient()
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
	// resp, err := client.PostForm("https://api.twitter.com/1.1/statuses/update.json",params)
	path := "https://api.twitter.com/1.1/statuses/update.json"
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
	// body, _ := ioutil.ReadAll(resp.Body)
	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	// err = json.Unmarshal(body, &responseTweet)
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
    // if len(load.TweetCreateEvent) < 1 || load.UserId == load.TweetCreateEvent[0].User.IdStr {
			if len(load.TweetCreateEvent) < 1 {
			// fmt.Println("Tweeted by bot...")
        return
    }
    //Send Hello world as a reply to the tweet, replies need to begin with the handles
    //of accounts they are replying to
    _, err = SendTweet("@"+load.TweetCreateEvent[0].User.Handle+" Hello World", load.TweetCreateEvent[0].IdStr)
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
	fmt.Println(string(responseJson))
	fmt.Fprintf(writer, string(responseJson))
}

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
