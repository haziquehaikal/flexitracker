package model

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

const (
	addProductStatement = `
	INSERT INTO product
	(proid,vendorID,productName,proinfo,typID,
	price,sprice,gprofit,issue,qty)
	VALUES(?,?,?,?,?,?,?,?,?,?)`

	getProDetailsStatement = `
	SELECT productName,sprice,proid FROM product`

	addProductOK  = `{"insert_product":true}`
	addProductErr = `{"insert_product":false}`
)

type Product struct {
	productName string  `json:"product_name,omitempty"`
	sprice      float32 `json:"sprice,omitempty"`
	productId   string  `json:"product_id,omitempty"`
}

func GetProductList(db *sql.DB) ([]byte, error) {
	var product Product

	aaa := []mapdata{}

	rows, err := db.Query(getProDetailsStatement)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&product.productName,
			&product.sprice,
			&product.productId)

		if err != nil {
			log.Print("ss ", err)
			return nil, err
		}
		list := mapdata{
			"pro_name": product.productName,
			"price":    product.sprice,
			"proid":    product.productId,
		}

		aaa = append(aaa, list)

	}

	defer rows.Close()
	fnl, err := json.Marshal(aaa)
	if err != nil {
		return nil, err
	}

	return fnl, err
}

func AddNewProduct(itemname string, itemprice, sellprice float32, db *sql.DB) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	log.Print(string(bytes))

	id := string(bytes)

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(addProductStatement)
	if err != nil {
		return "", err
	}
	profit := sellprice - itemprice
	_, err = stmt.Exec(id, "VEN001", itemname, "PHONE SPARTPART", "123", itemprice, sellprice, profit, "admin", 12)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}
	return id, nil

}

func DeleteProduct(id json.Number, db *sql.DB) ([]byte, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sqlstate := `SELECT * FROM product`
	rows, err := tx.Query(sqlstate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		//all attr
	}

	return nil, err
}

func ViewProduct(id json.Number, db *sql.DB) ([]byte, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sqlstate := `SELECT * FROM product`
	rows, err := tx.Query(sqlstate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		//all attr
	}

	return nil, err
}
