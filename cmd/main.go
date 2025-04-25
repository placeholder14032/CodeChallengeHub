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

    var err error
	_, err = database.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
    
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

            // // Admin credentials
            // username := "admin"
            // password := "admin123" // Change this to a secure password
        
            // // Create admin user
            // err = database.CreateAdminUser(username, password)
            // if err != nil {
            //     log.Fatalf("Failed to create admin user: %v", err)
            // }
            
    if err := app.Listen(); err != nil {
        log.Fatal(err)
    }

    // _, err := database.Connect()
	// if err != nil {
	// 	panic(err)
	// }


}
