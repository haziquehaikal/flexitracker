package cmd

import (
	"bitbucket.org/flexitracker/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var host string
var port string
var uname string
var pass string
var dbname string
var serveport string

var startCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the applicaions",
	Long:  `HAHAHAHAAH`,
	Run: func(cmd *cobra.Command, args []string) {
		//lang := cmd.Flag("start").Value.String()
		// listenPort := ":8080"
		// initalize an HTTP request multiplexer
		// mux := http.NewServeMux()

		// // Mount admin interface to mux
		// fs := http.FileServer(http.Dir("public"))
		// mux.Handle("/", fs)

		/*PRODUCTION USE*/

		host = viper.Get("DB_SERVICE_HOST").(string)
		//	port = viper.Get("DB_SERVICE_PORT").(string)
		uname = viper.Get("DB_SERVICE_USERNAME").(string)
		pass = viper.Get("DB_SERVICE_PASSWORD").(string)
		dbname = viper.Get("DB_SERVICE_DBNAME").(string)

		// host = "127.0.0.1"
		// uname = "root"
		// pass = ""
		// dbname = "flextech"
		router.InitApi(host, uname, pass, dbname)

		// fmt.Println("Listening on", listenPort)
		// http.ListenAndServe(listenPort, mux)
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
