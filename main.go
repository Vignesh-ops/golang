package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "vicky"
	dbname   = "postgres"
)

var err error

func connection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port = %d user= %s password =%s dbname = %s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	fmt.Println("The database is connected and error is ", err)
	return db
}

var tmpl *template.Template

type studentInfo struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

func register(w http.ResponseWriter, r *http.Request) {

	var tmplt = template.Must(template.ParseFiles("register.html"))
	tmplt.Execute(w, nil)
}

func insert(w http.ResponseWriter, r *http.Request) {
	db := connection()
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	student := studentInfo{
		Name:     r.FormValue("f_name"),
		Mail:     r.FormValue("mail"),
		Mobile:   r.FormValue("phone"),
		Password: r.FormValue("password"),
	}
	if r.FormValue("password") == r.FormValue("repassword") {
		_, err := db.Exec("insert into form(name,mail,mobile,password)values($1,$2,$3,$4)", student.Name, student.Mail, student.Mobile, student.Password)
		println("values are inserted", err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
	defer db.Close()

}
func getform(w http.ResponseWriter, r *http.Request) {
	db := connection()
	ss := []studentInfo{}
	s := studentInfo{}
	rows, err := db.Query("select * from form")
	fmt.Println(err)
	for rows.Next() {
		rows.Scan(&s.Name, &s.Mail, &s.Mobile, &s.Password)

		ss = append(ss, s)
	}

	json, _ := json.MarshalIndent(ss, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	defer rows.Close()
	defer db.Close()
}

func deleteform(w http.ResponseWriter, r *http.Request) {
	db := connection()
	_, err := db.Exec("delete  from form ")

	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "deleted")
	defer db.Close()
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/form", register).Methods("GET")
	router.HandleFunc("/list", getform).Methods("GET")
	router.HandleFunc("/delete", deleteform)
	router.HandleFunc("/myform", insert).Methods("POST")
	http.ListenAndServe(":8080", router)
}
