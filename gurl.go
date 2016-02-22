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
   Usage: ./gurl -site <websiteURL> */

func main() {

	// Get website URL
	flag.Parse()
	if site[0:7] != "http://" {
		site = "http://" + site
	}

	if check_url(site) == false {
		fmt.Println("Error: not a valid url.")
		return
	}

	// Hash the url
	hashed := md5.New()
	hashed.Write([]byte(site))
	hashed_bytes := hashed.Sum(nil)

	// base64 encode hashed url
	short_url := "localhost:8080/" + b64_encode(hashed_bytes)
	add_to_database(short_url, site)

	// Testing
	fmt.Println("Original url: ", site)
	fmt.Println("Shortened(?) url: ", short_url)

	// Create HTTP server
	http.HandleFunc("/", redirect)
	http.ListenAndServe(":8080", nil)
}

// Handle command line arguments
var site string

// Check url to see if it's valid
func check_url(site string) bool {
	// Create an HTTP client
	client := http.Client{}
	request, err := http.NewRequest("HEAD", site, nil)
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

func init() {
	flag.StringVar(&site, "site", "http://www.hollytancredi.net", "Website to connect to")
}

// Convert hashed bytes to base64
func b64_encode(hashed []byte) string {
	base64_encoded := base64.StdEncoding.EncodeToString(hashed)
	return base64_encoded[0:3]
}

// HTTP request handling
func redirect(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, url, http.StatusFound)
	fmt.Fprintf(w, "Test!")
}

// Map short url => original url in database
func add_to_database(short_url string, site string) {
	db, err := sql.Open("sqlite3", "./urls.db")
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
	}
	defer db.Close()

	create := `
        create table if not exists urls
        (id integer not null primary key,
        short text, original text);
        `
	_, err = db.Exec(create)
	if err != nil {
		fmt.Printf("%q: %s\n", err, create)
		return
	}

	// Add url map to database
	// tx => transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	insert, err := tx.Prepare("insert into urls(short, original) values(?,?)")
	if err != nil {
		fmt.Println("Error inserting into database: ", err)
	}
	defer insert.Close()

	_, err = insert.Exec(short_url, site)
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

	query := "select original from urls where short = " + short_url

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error querying database: ", err)
	}
	defer rows.Close()

	//For now, there should not be collisions
	var original_url string
	rows.Scan(&original_url)
	return original_url
}
