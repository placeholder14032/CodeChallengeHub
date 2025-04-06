package web

import (
    "fmt"
    "net/http"
    "time"

    "github.com/placeHolder143032/CodeChallengeHub/web/routes"

)

type Server struct {
    Host string
    PORT string
    URL string
}

func (a *Application) Listen() error {

    host := fmt.Sprintf("%s:%s", a.Server.Host, a.Server.PORT)

    server := http.Server{
        Handler:     a.routes(),
        Addr:        host,
        ReadTimeout: 300 * time.Second,
    }


    a.InfoLog.Printf("Server listening on :%s\n", host)

    return server.ListenAndServe()
}

func (a *Application) routes() http.Handler {
    mux := http.NewServeMux()

    // Frontend Routes (HTML pages)
    mux.HandleFunc("/", routes.LandingHandler)      // Landing page
    mux.HandleFunc("/auth", routes.AuthHandler)   // Signin-Login page

    mux.HandleFunc("/api/auth/login", routes.LoginHandler)    // User login
    mux.HandleFunc("/api/auth/register", routes.SigninHandler) // User registration


    return mux
}
