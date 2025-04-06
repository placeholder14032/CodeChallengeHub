package web

import (
    "fmt"
    "net/http"
    "time"
)

type Server struct {
    Host string
    PORT int
    URL string
}

func (a *Application) Listen() error {
    host := fmt.Sprintf("%s:%s", a.Server.Host, a.Server.PORT)

    srv := http.Server{
        // Handler:     a.routes(),
        Addr:        host,
        ReadTimeout: 300 * time.Second,
    }

    a.InfoLog.Printf("Server listening on :%s\n", host)

    return srv.ListenAndServe()
}
