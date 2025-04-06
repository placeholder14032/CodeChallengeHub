package web

import (
    "log"
)

type Application struct{
    AppName string
    Server *Server
    Debug bool
    ErrorLog *log.Logger
    InfoLog *log.Logger
}

func MyApp(appName string,server *Server,debug bool,errorLog *log.Logger,infoLog *log.Logger) *Application{
    return &Application{
        AppName: appName,
        Server: server,
        Debug: debug,
        ErrorLog: errorLog,
        InfoLog: infoLog,
    }
}