package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"url-app-backend/gen"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func welcomeFunction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

const (
	Db_USER     = "admin"
	Db_PASSWORD = "alypsok"
	Db_NAME     = "url"
)

func handleRequests() {
	http.HandleFunc("/post", postFunction)
	http.HandleFunc("/", getFunction)
	http.HandleFunc("/welcome", welcomeFunction)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {

	fmt.Println("It works.")

	DbInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s Dbname=%s sslmode=disable", "localhost", 5432, Db_USER, Db_PASSWORD, Db_NAME)
	// DbInfo := fmt.Sprintf("postgres://%s:%s@Db:5432/%s?sslmode=disable", Db_USER, Db_PASSWORD, Db_NAME)

	var err error
	Db, err = sql.Open("postgres", DbInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer Db.Close() 

	handleRequests()
}

type postUrl struct {
	Url string `json:"url"`
}

func getFunction(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path[1:]

	getQuery := fmt.Sprintf("SELECT long FROM link WHERE short='%s';", url)
	rows, err := Db.Query(getQuery)

	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer rows.Close()

	for rows.Next() {
		long := ""
		err := rows.Scan(&long)
		if err != nil {
			fmt.Println(err)
		} else {
			w.Write([]byte(fmt.Sprintf(`{"url": "%s"}`, long)))
		}
	}
}

func postFunction(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Post Read Error", http.StatusBadRequest)
	}

	var responseObject postUrl
	json.Unmarshal(body, &responseObject)

	url := responseObject.Url

	if url == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	getQuery := fmt.Sprintf("SELECT * FROM link WHERE short='%s';", url)
	rows, _ := Db.Query(getQuery)

	if !rows.Next() {
		short := gen.GetCode(Db)

		fmt.Println("3")
		insertStmt := fmt.Sprintf("INSERT INTO link (long, short) VALUES ('%s', '%s');", url, short)
		_, e := Db.Exec(insertStmt)

		if e != nil {
			w.Write([]byte(e.Error()))
		}

		w.Write([]byte(fmt.Sprintf(`{"url": "localhost/%s"}`, short)))
	} else {
		fmt.Println("Exception")
	}
}
