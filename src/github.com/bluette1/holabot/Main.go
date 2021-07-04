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

	//Create a new Mux Handler
	m := mux.NewRouter()
	//Listen to the base url and send a response
	m.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(200)
			fmt.Fprintf(writer, "Server is up and running")
	})
	//Listen to crc check and handle
	m.HandleFunc("/webhook/twitter", CrcCheck).Methods("GET")

	//Start Server
	server := &http.Server{
			Handler: m,
	}
	server.Addr = ":9090"
	server.ListenAndServe()
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
	// httpClient := CreateClient()
	var (
		err    error
		client *twittergo.Client
		req    *http.Request
		resp   *twittergo.APIResponse
		// user   *twittergo.User
	)
	client, err = LoadCredentials()
	if err != nil {
		fmt.Printf("Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}

	//Set parameters
	path := "/1.1/account_activity/all/"+os.Getenv("WEBHOOK_ENV")+"/webhooks.json"
  fmt.Println(path)
	values := url.Values{}
	appUrl := url.QueryEscape(os.Getenv("APP_URL")+"/webhook/twitter")
	values.Set("url", appUrl)

	body := strings.NewReader(values.Encode())
	req, err = http.NewRequest("POST", path, body)
	defer req.Body.Close()
	
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}

	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}

	// Parse response and check response
	respBody, err := ioutil.ReadAll(resp.Body)

    fmt.Println(string(respBody))
}
