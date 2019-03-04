package model

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"
)

type mapdata map[string]interface{}

const (
	host   = "localhost"
	port   = 5432
	uname  = "postgres"
	passwd = "tuxuri"
	dbname = "tuxuri"
)

//host, uname, pass string, dbname string
func RunDB() (*sql.DB, error) {

	viper.GetString("TEST")
	// info := fmt.Sprintf("host=%s port=%s user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, uname, passwd, dbname)
	db, err := sql.Open("mysql", "root:wolwolf96@/flextech")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, err

}

func CheckDB(db *sql.DB) error {
	sqlstate := "SELECT EXISTS(SELECT 1 from information_schema.tables WHERE table_schema = 'public' AND table_name = 'geodata');"
	rows, err := db.Query(sqlstate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var exists string
		err := rows.Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(exists)
		if exists == "false" {
			crttbl := "CREATE TABLE GEODATA(gid VARCHAR(50) PRIMARY KEY ,latitude float8 NOT NULL, longitude float8 NOT NULL, zoom float8 NOT NULL,flag int NOT NULL);"
			_, err := db.Exec(crttbl)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Table successfully created")
		} else {
			log.Println("You are good to go :)")
		}
	}
	return err
}
