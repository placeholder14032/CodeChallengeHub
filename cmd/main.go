package main

import (
	"fmt"
	"log"
	"os"

	"github.com/placeHolder143032/CodeChallengeHub/web"
    "github.com/placeHolder143032/CodeChallengeHub/database"

)

func main() {
    host := "localhost"
    port := "8090"
    
    // Initialize the application
    app := &web.Application{
        AppName: "CodeChallengeHub",
        Server: &web.Server{
            Host: host,
            PORT: port,
            URL: fmt.Sprintf("http://%s:%s", host, port),
        },
        Debug: true,
        ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
        InfoLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile),
    }

    if err := app.Listen(); err != nil {
        log.Fatal(err)
    }

    _, err := database.Connect()
	if err != nil {
		panic(err)
	}
}
