package main

import (
    "fmt"
    "log"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "net/http"
)

func main(){
    //Load env
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
        fmt.Println("Error loading .env file")
    }
    fmt.Println("Starting Server")

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
    //TODO implement CRC check
}
