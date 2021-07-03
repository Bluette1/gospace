package main

import (
    "os"
    "fmt"
    "bytes"
    "log"
    "io/ioutil"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/dghubble/oauth1"
    "net/http"
    "net/url"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
)

func main(){
    //Load env
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
    fmt.Fprintf(writer, string(responseJson))
}

func CreateClient() *http.Client {
    //Create oauth client with consumer keys and access token
    config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
    token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))
    // subscribeWebhook()

    return config.Client(oauth1.NoContext, token)
}
func registerWebhook(){
    fmt.Println("Registering webhook...")
    httpClient := CreateClient()

    //Set parameters
    path := "https://api.twitter.com/1.1/account_activity/webhooks.json"

    values := url.Values{}
    values.Set("url", os.Getenv("APP_URL")+"/webhook/twitter")

    //Make Oauth Post with parameters
    resp, _ := httpClient.PostForm(path, values)
    defer resp.Body.Close()
    //Parse response and check response
    body, _ := ioutil.ReadAll(resp.Body)
    var data map[string]interface{}
    if err := json.Unmarshal([]byte(body), &data); err != nil {
        panic(err)
    }
    fmt.Println("Webhook id of " + data["id"].(string) + " has been registered")
}

func subscribeWebhook(){
    fmt.Println("Subscribing webapp...")
    client := CreateClient()
    path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/subscriptions.json"
    resp, _ := client.PostForm(path, nil)
    body, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    //If response code is 204 it was successful
    if resp.StatusCode == 204 {
        fmt.Println("Subscribed successfully")
    } else if resp.StatusCode!= 204 {
        fmt.Println("Could not subscribe the webhook. Response below:")
        fmt.Println(string(body))
    }
}