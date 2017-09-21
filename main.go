package main

import (
	appConfig "VCMS/apps/config"
	"VCMS/apps/route"

	"github.com/Sirupsen/logrus"
	 
)

func init() {
	// config Loading
	appConfig.LoadAutoConfig()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	router := route.Init()

	//router.StartServer()
	router.Logger.Fatal(router.Start(appConfig.Config.APP.Port))

}
