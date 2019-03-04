package main

import (
	"fmt"

	"bitbucket.org/flexitracker/router"
)

var host string
var port string
var uname string
var pass string
var dbname string
var serveport string

func main() {

	// viper.AutomaticEnv()
	fmt.Println("Welcome To Flexi")
	// host = viper.Get("DB_SERVICE_HOST").(string)
	// log.Println(host)
	// uname = viper.Get("DB_SERVICE_USERNAME").(string)
	// pass = viper.Get("DB_SERVICE_PASSWORD").(string)
	// dbname = viper.Get("DB_SERVICE_DBNAME").(string)
	//host, uname, pass, dbname
	router.InitApi()
	//cmd.Execute()

}
