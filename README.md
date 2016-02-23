# :sparkles: gurl :sparkles:
### gurl is a url shortener written in go.

Since I just wrote this for fun, storage is an `sqlite
database` and the server is `localhost` listening to port `8080`.

### Usage
You can use the optional -url flag to specify the url you want to shorten.
`./gurl -url www.github.com`   
Otherwise, the server just redirects based on previously shortened urls
found in the database.

## Dependencies
* sqlite3
* mattn's [go-sqlite3](https://github.com/mattn/go-sqlite3) package
   
   
   
   
   
:raising_hand:
