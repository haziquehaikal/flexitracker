package router

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"bitbucket.org/flexitracker/model"
)

const (

	//Error Responses
	httpOK          = `{"response":"OK"}`
	httpErrMthd     = `{"response":"Method not allowed"}`
	httpInternalErr = `{"response":"Internal server Error"}`
)

type User struct {
	Username string `json:"uname"`
	Password string `json:"passwd"`
}

type Job struct {
	CustName    string   `json:"cust_name,omitempty"`
	CustEmail   string   `json:"cust_email,omitempty"`
	CustNo      string   `json:"cust_no,omitempty"`
	ItemUsed    []string `json:"item_used,omitempty"`
	TotalAmount float32  `json:"total_amount,omitempty"`
	StaffId     string   `json:"staff_id,omitempty"`
	ProId       string   `json:"pro_id,omitempty"`
	Qty         int      `json:"qty,omitempty"`
	AddMode     string   `json:"add_mode,omitempty"`
	ItemName    string   `json:"item_name,omitempty"`
	ItemPrice   float32  `json:"item_price,omitempty"`
	SellPrice   float32  `json:"sell_price,omitempty"`
}

type Staff struct {
	StaffId string `json:"staff_id,omitempty"`
}

type StaffRegister struct {
	StaffEmail    string `json:"staff_email,omitempty"`
	StaffPassword string `json:"staff_password,omitempty"`
	StaffId       string `json:"staff_id,omitempty"`
}

type App struct {
	db *sql.DB
}

//host string, uname string, pass string, dbname string
func InitApi() {

	tx, err := model.RunDB()
	if err != nil {
		log.Fatal(err)
	}
	db := &App{db: tx}

	checkapi := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("kucintankandikau"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	serveport := "5000"
	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("./ui/dist/")))
	/*Auth endpoint*/
	router.HandleFunc("/flexitracker/api/v1/login", db.loginHandler)
	router.HandleFunc("/flexitracker/api/v1/register", db.userRegisterHandler)
	router.HandleFunc("/flexitracker/api/v1/logout", db.loginHandler)
	//router.HandleFunc("/flexitracker/api/v1/job/insight/{staffid}", db.jobInsightHandler)
	/*JOB endpoint*/
	router.Handle("/flexitracker/api/v1/job/list/{staffid}", checkapi.Handler(http.HandlerFunc(db.jobListHandler)))
	router.Handle("/flexitracker/api/v1/job/insight/{staffid}", checkapi.Handler(http.HandlerFunc(db.jobInsightHandler)))
	router.Handle("/flexitracker/api/v1/job/add", checkapi.Handler(http.HandlerFunc(db.jobAddHandler)))
	router.Handle("/flexitracker/api/v1/job/add/custom", http.HandlerFunc(db.jobAddHandler))
	router.Handle("/flexitracker/api/v1/job/view", checkapi.Handler(http.HandlerFunc(db.jobListHandler)))
	/*PRODUCT endpoint*/
	router.Handle("/flexitracker/api/v1/product/list", checkapi.Handler(http.HandlerFunc(db.productListHandler)))
	router.Handle("/flexitracker/api/v1/product/add", checkapi.Handler(http.HandlerFunc(db.jobListHandler)))
	err = http.ListenAndServe(":"+serveport, router)
	if err != nil {
		log.Fatal(err)
	}

}

func (a *App) IndexHandler(w http.ResponseWriter, r *http.Request) {

}

func (a *App) insertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "	*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

}

/*AUTH*/
func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var data User

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Print(err)
			http.Error(w, httpInternalErr, http.StatusInternalServerError)
		}

		uname := data.Username
		passwd := data.Password

		//log.Print(uname, passwd)

		v, err := model.CheckLogin(uname, passwd, a.db)
		if err != nil {
			log.Print(err)
		}
		w.Write(v)

	} else {
		http.Error(w, "INVALID METHOD ", http.StatusMethodNotAllowed)
	}
}

func (a *App) userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var data StaffRegister
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Print(err)
			http.Error(w, httpInternalErr, http.StatusInternalServerError)
		}

		res, err := model.AddNewUser(data.StaffEmail, data.StaffPassword, data.StaffId, a.db)
		if err != nil {
			http.Error(w, httpInternalErr, http.StatusInternalServerError)
			log.Printf("User registration failed ,Message: %v", err)
		}
		w.Write(res)

	} else {
		http.Error(w, httpErrMthd, http.StatusMethodNotAllowed)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	//kill session
}

/*JOB*/

func (a *App) jobAddHandler(w http.ResponseWriter, r *http.Request) {
	var jobdata Job
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &jobdata)
		if err != nil {
			log.Printf("Cannot unmarshal data, Message: %v", err)
		}

		if jobdata.AddMode == "custom" {
			proid, err := model.AddNewProduct(jobdata.ItemName, jobdata.ItemPrice, jobdata.SellPrice, a.db)
			if err != nil {
				log.Printf("Cannot add product , Message: %v", err)
			}
			res, err := model.SaveJobDone(jobdata.CustName, jobdata.CustEmail, jobdata.CustNo, jobdata.TotalAmount, jobdata.StaffId, jobdata.Qty, proid, a.db)
			if err != nil {
				log.Printf("cannot add custom job ,Message: %v", err)
				panic(err)

			}
			//log.Print(tet)
			w.Write(res)

		}
		if jobdata.AddMode == "normal" {
			res, err := model.SaveJobDone(jobdata.CustName, jobdata.CustEmail, jobdata.CustNo, jobdata.TotalAmount, jobdata.StaffId, jobdata.Qty, jobdata.ProId, a.db)
			if err != nil {
				panic(err)
			}

			w.Write(res)
		}

	} else {
		http.Error(w, "INVALID METHOD", http.StatusMethodNotAllowed)
	}

}

func (a *App) jobListHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		obj := mux.Vars(r)
		res, err := model.GetPreviousJob(obj["staffid"], obj["type"], a.db)
		if err != nil {
			log.Print(err)
			http.Error(w, httpInternalErr, http.StatusInternalServerError)
		}
		w.Write(res)
	} else {
		http.Error(w, "INVALID METHOD", http.StatusMethodNotAllowed)
	}

}

func (a *App) jobInsightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		obj := mux.Vars(r)
		res, err := model.JobStatistic(obj["staffid"], a.db)
		if err != nil {
			log.Print(err)
		}

		w.Write(res)
	}
}

/* PRODUCT */

func (a *App) productListHandler(w http.ResponseWriter, r *http.Request) {
	v, err := model.GetProductList(a.db)
	if err != nil {
		log.Print(err)
		http.Error(w, httpInternalErr, http.StatusInternalServerError)
	}
	w.Write(v)
}

func mobileLogin(w http.ResponseWriter, r *http.Request) {

	// idGen := strconv.Itoa((rand.Intn(9999)))
	// var randstr = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	// fix := make([]byte, 5)
	// var s string
	// for i := 0; i < 5; i++ {
	// 	fix[i] = randstr[rand.Intn(len(randstr))]
	// 	s += string(fix[i])
	// }
	// finalID := s + idGen
	// fmt.Println(finalID)
}
