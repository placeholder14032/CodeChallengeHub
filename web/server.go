package web

type Server struct {
    host string
    port int
    url string
}

func (a *Application) listenAndServer() error {
    host := fmt.Sprintf("%s:%s", a.server.host, a.server.port)

    srv := http.Server{
        Handler:     a.routes(),
        Addr:        host,
        ReadTimeout: 300 * time.Second,
    }

    a.infoLog.Printf("Server listening on :%s\n", host)

    return srv.ListenAndServe()
}