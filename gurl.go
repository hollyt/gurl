package main

import (
    "crypto/md5"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
)

// gurl: a URL shortener written in Go
// Usage: ./gurl -site <websiteURL>

var site string
func init() {
    flag.StringVar(&site, "site", "http://www.hollytancredi.net", "Website to connect to")
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
    hashed_string := hashed.Sum(nil)

    // Testing
    fmt.Println("Original url: ", site)
    fmt.Printf("Hashed url: %x\n", hashed_string)
}
