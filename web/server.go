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
    mux.HandleFunc("/", routes.TempHandler)
    return mux
}
