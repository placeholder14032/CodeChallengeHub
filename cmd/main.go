package main

import (
    "fmt"
    "log"
    "os"

    "github.com/placeHolder143032/CodeChallengeHub/web"
)

func main(){

    // fmt.Println("Hello, World!")
    host := "localhost"
    port := 8080
    fmt.Printf("Server is running on %s:%d\n", host, port)
    // Initialize the application
    app := &web.Application{
        AppName: "CodeChallengeHub",
        Server: &web.Server{
            Host: host,
            PORT: port,
        },
        Debug: true, // we can have auto run server for debug mode
        ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
        InfoLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile),
    }
    app.Listen()
    // time.Sleep(5 * time.Second)
}