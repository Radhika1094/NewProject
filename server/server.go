package main

import (
	"CGScrutinyForm/server/api"
	"CGScrutinyForm/server/api/model"

	"log"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"
	"github.com/labstack/echo"
)

func main() {
	confighelper.InitViper()
	e := echo.New()
	e.HideBanner = true
	//should be called before initializing api
	initializeVariable()
	//Bind API
	api.Init(e)
	serverPort := confighelper.GetConfig("serverport")
	logginghelper.LogDebug(serverPort)
	err := e.Start(serverPort)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeVariable() {
	model.DBNAME = confighelper.GetConfig("DBNAME")
}
