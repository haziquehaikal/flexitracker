package model

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

const (
	getJobStatement = `SELECT custname ,custemail 
	, custno,totalamount,
	job_id,date_created ,productName
	FROM jobs INNER JOIN product
	WHERE jobs.proid = product.proid
	AND staff_id = ? ORDER BY 
	date_created DESC`

	addJobStatement = `INSERT INTO jobs
	 (job_id,custname ,custemail , 
	custno,totalamount,staff_id,qty,proid) 
	 VALUES (?,?,?,?,?,?,?,?)`

	getJobStatistic = `select 
	 SUM(totalamount) 
	 from jobs 
	 WHERE staff_id = ?`
)

type Job struct {
	CustName          sql.NullString  `json:"cust_name,omitempty"`
	CustEmail         sql.NullString  `json:"cust_email,omitempty"`
	CustNo            sql.NullString  `json:"cust_no,omitempty"`
	ItemUsed          sql.NullString  `json:"item_used,omitempty"`
	TotalAmount       sql.NullFloat64 `json:"total_amount,omitempty"`
	JobDate           sql.NullString  `json:"job_date,omitempty"`
	JobId             sql.NullString  `json:"job_id,omitempty"`
	ProductName       sql.NullString
	TotalStaffCollect sql.NullFloat64 `json:"total_staff_collect,omitempty"`
}

func GetPreviousJob(a string, b string, db *sql.DB) ([]byte, error) {
	var jobdata Job
	d := []mapdata{}
	rows, err := db.Query(getJobStatement, a)
	if err != nil && err == sql.ErrNoRows {
		//log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&jobdata.CustName, &jobdata.CustEmail,
			&jobdata.CustNo, &jobdata.TotalAmount,
			&jobdata.JobId, &jobdata.JobDate, &jobdata.ProductName)
		if err != nil {
			return nil, err
		}

		joblist := mapdata{
			"cust_name":    jobdata.CustName.String,
			"cust_email":   jobdata.CustEmail.String,
			"cust_no":      jobdata.CustNo.String,
			"total_amount": jobdata.TotalAmount.Float64,
			"job_date":     jobdata.JobDate.String,
			"job_id":       jobdata.JobId.String,
			"product_name": jobdata.ProductName.String,
		}
		d = append(d, joblist)
	}

	final, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return final, err
}

func JobStatistic(staffid string, db *sql.DB) ([]byte, error) {
	var jobdata Job

	stmt, err := db.Prepare(getJobStatistic)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(staffid).Scan(&jobdata.TotalStaffCollect)
	if err != nil {
		return nil, err
	}

	res := mapdata{
		"total_collect": jobdata.TotalStaffCollect.Float64,
	}

	final, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return final, err
}

func SaveJobDone(custname string, custemail string, custphone string, totalamount float32, staffid string, qty int, proid string, db *sql.DB) ([]byte, error) {

	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	log.Print(string(bytes))
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(addJobStatement)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(string(bytes), custname, custemail, custphone, totalamount, staffid, qty, proid)
	if err != nil {
		log.Print(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Print(err)
	}

	if err != nil {
		res := mapdata{"JOB_INSERT": false}
		v, _ := json.Marshal(res)
		return v, err
	} else {
		res := mapdata{"JOB_INSERT": true}
		v, _ := json.Marshal(res)
		return v, nil
	}

}

func SaveCustomJob(custname string, custemail string, custphone string, totalamount float32, staffid string, qty int, proid string, db *sql.DB) ([]byte, error) {

	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	log.Print(string(bytes))
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(addJobStatement)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(string(bytes), custname, custemail, custphone, totalamount, staffid, qty, proid)
	if err != nil {
		log.Print(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Print(err)
	}

	if err != nil {
		res := mapdata{"JOB_INSERT": false}
		v, _ := json.Marshal(res)
		return v, err
	} else {
		res := mapdata{"JOB_INSERT": true}
		v, _ := json.Marshal(res)
		return v, nil
	}

}
