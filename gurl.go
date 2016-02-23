package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

/* gurl: a URL shortener written in Go
   Usage: ./gurl -url <websiteURL> */

func main() {
	flag.Parse()
	fmt.Println("url: ", url)
	if url != " " {
		shorten()
	}

	// Create HTTP server
	http.HandleFunc("/", redirect)
	http.ListenAndServe(":8080", nil)
}

// Handle command line arguments
var url string
var short_url string

// Check url to see if it's valid
func check_url(url string) bool {
	// Create an HTTP client
	client := http.Client{}
	request, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request: ", err)
		return false
	}

	// Send the request and get back the HTTP response
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending HTTP request: ", err)
		return false
	}
	defer response.Body.Close()
	return true
}

// Shorten url if one is entered
func shorten() {
	// Get website URL
	if url[0:7] != "http://" {
		url = "http://" + url
	}

	if check_url(url) == false {
		fmt.Println("Error: not a valid url.")
		return
	}

	// Hash the url
	hashed := md5.New()
	hashed.Write([]byte(url))
	hashed_bytes := hashed.Sum(nil)

	// base64 encode hashed url
	short_url = "localhost:8080/" + b64_encode(hashed_bytes)
	add_to_database(short_url, url)

	// Testing
	fmt.Println("Original url: ", url)
	fmt.Println("Shortened(?) url: ", short_url)
}

func init() {
	flag.StringVar(&url, "url", " ", "URL to shorten")
}

// Convert hashed bytes to base64
func b64_encode(hashed []byte) string {
	base64_encoded := base64.StdEncoding.EncodeToString(hashed)
	return base64_encoded[0:3]
}

// HTTP request handling
func redirect(w http.ResponseWriter, r *http.Request) {
	redirect_url := get_original_url("localhost:8080" + r.URL.Path)
	http.Redirect(w, r, redirect_url, 302)
	//fmt.Fprintf(w, redirect_url)
}

// Map short url => original url in database
func add_to_database(short_url string, url string) {
	db, err := sql.Open("sqlite3", "./urls.db")
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
	}
	defer db.Close()

	create := `
        create table if not exists urls
        (short text primary key, original text);
        `
	_, err = db.Exec(create)
	if err != nil {
		fmt.Printf("%q: %s\n", err, create)
		return
	}

	// Add url map to database
	// tx is transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	insert, err := tx.Prepare("insert or ignore into urls(short, original) values(?,?)")
	if err != nil {
		fmt.Println("Error inserting into database: ", err)
	}
	defer insert.Close()

	_, err = insert.Exec(short_url, url)
	if err != nil {
		fmt.Println("Error inserting into database: ", err)
	}
	tx.Commit()
}

func get_original_url(short_url string) string {
	db, err := sql.Open("sqlite3", "./urls.db")
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
	}
	defer db.Close()

	var original_url string
	query, err := db.Prepare("select original from urls where short = ?")
	if err != nil {
		fmt.Println("Error preparing query: ", err)
	}
	err = query.QueryRow(short_url).Scan(&original_url)
	if err != nil {
		fmt.Println("Error querying database: ", err)
	}

	return original_url
}
