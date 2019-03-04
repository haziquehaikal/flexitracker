package model

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	//SQL query

	checkUserStatement = `SELECT staff_id,staff_pass,
	staff_name FROM staff WHERE email = ?`

	createTableStatement = `CREATE TABLE IF NOT EXISTS GEODATA(
		geodata_id serial PRIMARY KEY
		,latitude float8 NOT NULL, 
		longitude float8 NOT NULL,
		zoom float8 NOT NULL,
		flag int NOT NULL)`

	insertUserStatement = `INSERT INTO staff VALUES (?,?,?,?)`

	viewStatement = `SELECT 
	geodata_id, x, y, zoom, flag, ST_AsGeoJSON(
	ST_Transform(TileBBox(zoom, x, y), 4326), 6) 
	  FROM geodata`

	deleteStatement = `DELETE FROM geodata WHERE geodata_id = $1`

	updateStatement = `UPDATE geodata SET flag  = $1 WHERE geodata_id = $2`

	//Error Responses

	registerOK = `{"DATA DELETE":true}`
	deleteErr  = `{"DATA DELETE":false}`
	updateOK   = `{"DATA UPDATE":true}`
	updateErr  = `{"DATA UPDATE":false}`
	viewOK     = `{"DATA VIEW":true}`
	viewErr    = `{"DATA VIEW":false}`
)

type DB struct {
	*sql.DB
}

type User struct {
	staffName  sql.NullString `json:"staff_name,omitempty"`
	staffId    sql.NullString `json:"staff_id,omitempty"`
	staffEmail sql.NullString `json:"staff_email,omitempty"`
	earnId     sql.NullString `json:"earn_id,omitempty"`
	staffPass  sql.NullString `json:"staff_pass,omitempty"`
	totalEarn  sql.NullString `json:"total_earn,omitempty"`
}

var key = []byte("kucintankandikau")

// CheckLogin i know this is a stupid login function , SHUT UP !
func CheckLogin(email string, password string, db *sql.DB) ([]byte, error) {
	var user User
	msg := mapdata{}
	log.Printf("Username: %v , Password: %v", email, password)

	err := db.QueryRow(checkUserStatement, email).Scan(&user.staffId, &user.staffPass, &user.staffName)
	if err != nil && err != sql.ErrNoRows {

		log.Printf("Cannot find email, Message: %v", err)
		msg = mapdata{
			"LOGIN_VERIFY": false,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.staffPass.String), []byte(password))
	if err != nil {
		msg = mapdata{
			"LOGIN_VERIFY": false,
		}
	} else {

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)

		claims["email"] = user.staffEmail
		claims["name"] = user.staffName
		claims["id"] = user.staffId
		claims["exp"] = time.Now().Add(time.Hour * 999999).Unix()

		tokenString, _ := token.SignedString(key)

		msg = mapdata{
			"LOGIN_VERIFY": true,
			"STAFF_EMAIL":  email,
			"STAFF_NAME":   user.staffName.String,
			"STAFF_ID":     user.staffId.String,
			"JWT_TOKEN":    tokenString,
		}

	}

	v, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return v, err

}

func GetUserDetails(db *sql.DB) ([]byte, error) {
	var user User
	sqlstate := `SELECT staff_email , staff_name FROM staff`
	rows, err := db.Query(sqlstate)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.staffEmail, &user.staffName)
		if err != nil {
			return nil, err
		}
	}
	arr := []string{user.staffEmail.String, user.staffName.String}
	data, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}

	return data, err

}

func AddNewUser(username, password, staffid string, db *sql.DB) ([]byte, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(insertUserStatement)
	if err != nil {
		return nil, err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Print(err)
	}

	log.Printf("Username: %v , Password: %v", username, password)

	_, err = stmt.Exec(staffid, bytes, username, "bitch")
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	stmt.Close()

	if err != nil {
		res := mapdata{"register_status": false}
		v, _ := json.Marshal(res)
		return v, err
	} else {
		res := mapdata{"register_status": true}
		v, _ := json.Marshal(res)
		return v, nil
	}
}
