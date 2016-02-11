package main

import (
    "crypto/md5"
    "encoding/base64"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
)

/* gurl: a URL shortener written in Go
   Usage: ./gurl -site <websiteURL> */

// Handle command line arguments
var site string
func init() {
    flag.StringVar(&site, "site", "http://www.hollytancredi.net", "Website to connect to")
}


// Convert hashed bytes to base64
func b64_encode(hashed []byte) string {
    base64_encoded := base64.StdEncoding.EncodeToString(hashed)
    return base64_encoded
}

func main() {

    // Get website URL
    flag.Parse()
    if site[0:7] != "http://" {
        site = "http://" + site
    }

    // Create an HTTP client
    client := http.Client{}
    request, err := http.NewRequest("GET", site, nil)
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

    responseBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("Error getting response body: ", err)
        return
    }
    fmt.Println(string(responseBytes))

    // Hash the url
    hashed := md5.New()
    hashed.Write([]byte(site))
    hashed_bytes := hashed.Sum(nil)

    // base64 encode hashed url
    b64_url := b64_encode(hashed_bytes)

    // Testing
    fmt.Println("Original url: ", site)
    fmt.Printf("Hashed url (hex): %x\n", hashed_bytes)
    fmt.Println("Hashed url (bytes): ", hashed_bytes)
    fmt.Println("Hashed url (base64) :", b64_url)
}