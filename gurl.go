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

	// Create an HTTP client
	client := http.Client{}
	request, err := http.NewRequest("HEAD", site, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request: ", err)
		return
	}

	// Send the request and get back the HTTP response
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending HTTP request: ", err)
		return
	}
	defer response.Body.Close()

	// Hash the url
	hashed := md5.New()
	hashed.Write([]byte(site))
	hashed_bytes := hashed.Sum(nil)

	// base64 encode hashed url
	short_url := "h0lt.net/" + b64_encode(hashed_bytes)

        // Add short url => original url to database
        db, err := sql.Open("sqlite3", "./urls.db")
        if err != nil {
            fmt.Println("Error connecting to database: ", err)
        }
        defer db.Close()

        create := `
        create table test (id integer not null primary key, shortened text, original text);
        `
        _, err = db.Exec(create)
        if err != nil {
            fmt.Printf("%q: %s\n", err, create)
            return
        }

	// Testing
	fmt.Println("Original url: ", site)
	fmt.Println("Shortened(?) url: ", short_url)

	// Create HTTP server
	http.HandleFunc("/", redirect)
	http.ListenAndServe(":8080", nil)
}

// Handle command line arguments
var site string

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
	//short_url := "h0lt.net" + r.URL.Path
	//msg := short_url + " mapped to " + url_map[short_url]
	fmt.Fprintf(w, "Test!")
}
