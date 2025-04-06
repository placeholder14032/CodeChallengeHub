package main

import (
    "log"
	"github.com/placeHolder143032/CodeChallengeHub/web"
)

type Application struct{
    appName string
    server *Server
    deybug bool
    errorLog *log.Logger
    infoLog *log.Logger
}